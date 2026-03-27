package dao

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func InitSQLite(ctx context.Context, dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	for _, statement := range []string{
		`PRAGMA busy_timeout = 5000`,
		`PRAGMA journal_mode = WAL`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			_ = db.Close()
			return nil, err
		}
	}

	if err := bootstrapSQLite(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func bootstrapSQLite(ctx context.Context, db *sql.DB) error {
	schemaStatements := []string{
		`CREATE TABLE IF NOT EXISTS site_profile (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			name TEXT NOT NULL,
			headline TEXT NOT NULL,
			intro TEXT NOT NULL,
			location TEXT NOT NULL,
			domain TEXT NOT NULL,
			email TEXT NOT NULL,
			motto TEXT NOT NULL,
			tech_stack_json TEXT NOT NULL,
			stats_json TEXT NOT NULL,
			social_links_json TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT NOT NULL UNIQUE,
			title TEXT NOT NULL,
			summary TEXT NOT NULL,
			category TEXT NOT NULL,
			read_time TEXT NOT NULL,
			cover_label TEXT NOT NULL,
			content_markdown TEXT NOT NULL DEFAULT '',
			tags_json TEXT NOT NULL,
			featured INTEGER NOT NULL DEFAULT 0,
			published_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS post_assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT NOT NULL,
			mime_type TEXT NOT NULL,
			data_blob BLOB NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			summary TEXT NOT NULL,
			status TEXT NOT NULL,
			link TEXT NOT NULL,
			image_url TEXT NOT NULL DEFAULT '',
			accent TEXT NOT NULL,
			tech_stack_json TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS videos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			url TEXT NOT NULL,
			thumbnail_url TEXT NOT NULL DEFAULT '',
			published_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS timeline_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			period TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			author_name TEXT NOT NULL,
			content TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'approved',
			created_at TEXT NOT NULL,
			FOREIGN KEY(post_id) REFERENCES posts(id)
		);`,
		`CREATE TABLE IF NOT EXISTS post_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			visitor_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			UNIQUE(post_id, visitor_id),
			FOREIGN KEY(post_id) REFERENCES posts(id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);`,
		`CREATE INDEX IF NOT EXISTS idx_post_likes_post_id ON post_likes(post_id);`,
	}

	for _, statement := range schemaStatements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}

	if _, err := db.ExecContext(ctx, `ALTER TABLE videos ADD COLUMN thumbnail_url TEXT NOT NULL DEFAULT ''`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE projects ADD COLUMN image_url TEXT NOT NULL DEFAULT ''`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return err
	}
	_, err := migratePostMarkdownSchema(ctx, db)
	if err != nil {
		return err
	}

	return seedSQLite(ctx, db)
}

func seedSQLite(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := seedSiteProfileRow(ctx, tx); err != nil {
		return err
	}
	if err := seedPostRows(ctx, tx); err != nil {
		return err
	}
	if err := seedProjectRows(ctx, tx); err != nil {
		return err
	}
	if err := seedVideoRows(ctx, tx); err != nil {
		return err
	}
	if err := seedTimelineRows(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

func seedSiteProfileRow(ctx context.Context, tx *sql.Tx) error {
	site := seedSiteProfile()
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO site_profile
			(id, name, headline, intro, location, domain, email, motto, tech_stack_json, stats_json, social_links_json)
		VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			headline = excluded.headline,
			intro = excluded.intro,
			location = excluded.location,
			domain = excluded.domain,
			email = excluded.email,
			motto = excluded.motto,
			tech_stack_json = excluded.tech_stack_json,
			stats_json = excluded.stats_json,
			social_links_json = excluded.social_links_json`,
		site.Name,
		site.Headline,
		site.Intro,
		site.Location,
		site.Domain,
		site.Email,
		site.Motto,
		mustJSON(site.TechStack),
		mustJSON(site.Stats),
		mustJSON(site.SocialLinks),
	)
	return err
}

func seedPostRows(ctx context.Context, tx *sql.Tx) error {
	if hasRows, err := hasAnyRows(ctx, tx, "SELECT COUNT(1) FROM posts"); err != nil {
		return err
	} else if hasRows {
		return nil
	}

	for _, post := range seedPosts() {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO posts
				(slug, title, summary, category, read_time, cover_label, content_markdown, tags_json, featured, published_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			post.Slug,
			post.Title,
			post.Summary,
			post.Category,
			post.ReadTime,
			post.CoverLabel,
			post.ContentMarkdown,
			mustJSON(post.Tags),
			boolToInt(post.Featured),
			post.PublishedAt.Format(timeLayout),
		); err != nil {
			return err
		}
	}
	return nil
}

func migratePostMarkdownSchema(ctx context.Context, db *sql.DB) (bool, error) {
	hasMarkdownColumn, err := tableColumnExists(ctx, db, "posts", "content_markdown")
	if err != nil {
		return false, err
	}
	hasHeroNoteColumn, err := tableColumnExists(ctx, db, "posts", "hero_note")
	if err != nil {
		return false, err
	}
	hasBlocksColumn, err := tableColumnExists(ctx, db, "posts", "blocks_json")
	if err != nil {
		return false, err
	}
	if hasMarkdownColumn && !hasHeroNoteColumn && !hasBlocksColumn {
		return false, nil
	}

	for _, statement := range []string{
		`DROP TABLE IF EXISTS comments`,
		`DROP TABLE IF EXISTS post_likes`,
		`DROP TABLE IF EXISTS posts`,
		`DELETE FROM post_assets`,
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			slug TEXT NOT NULL UNIQUE,
			title TEXT NOT NULL,
			summary TEXT NOT NULL,
			category TEXT NOT NULL,
			read_time TEXT NOT NULL,
			cover_label TEXT NOT NULL,
			content_markdown TEXT NOT NULL DEFAULT '',
			tags_json TEXT NOT NULL,
			featured INTEGER NOT NULL DEFAULT 0,
			published_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			author_name TEXT NOT NULL,
			content TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'approved',
			created_at TEXT NOT NULL,
			FOREIGN KEY(post_id) REFERENCES posts(id)
		)`,
		`CREATE TABLE IF NOT EXISTS post_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			visitor_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			UNIQUE(post_id, visitor_id),
			FOREIGN KEY(post_id) REFERENCES posts(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id)`,
		`CREATE INDEX IF NOT EXISTS idx_post_likes_post_id ON post_likes(post_id)`,
	} {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return false, err
		}
	}

	return true, nil
}

func tableColumnExists(ctx context.Context, db *sql.DB, tableName, columnName string) (bool, error) {
	rows, err := db.QueryContext(ctx, `PRAGMA table_info(`+tableName+`)`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			columnType string
			notNull    int
			defaultVal sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultVal, &pk); err != nil {
			return false, err
		}
		if strings.EqualFold(name, columnName) {
			return true, nil
		}
	}

	return false, rows.Err()
}

func seedProjectRows(ctx context.Context, tx *sql.Tx) error {
	if hasRows, err := hasAnyRows(ctx, tx, "SELECT COUNT(1) FROM projects"); err != nil {
		return err
	} else if hasRows {
		return nil
	}

	for _, project := range seedProjects() {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO projects (name, summary, status, link, image_url, accent, tech_stack_json)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			project.Name,
			project.Summary,
			project.Status,
			project.Link,
			project.ImageURL,
			project.Accent,
			mustJSON(project.TechStack),
		); err != nil {
			return err
		}
	}
	return nil
}

func seedTimelineRows(ctx context.Context, tx *sql.Tx) error {
	if hasRows, err := hasAnyRows(ctx, tx, "SELECT COUNT(1) FROM timeline_entries"); err != nil {
		return err
	} else if hasRows {
		return nil
	}

	for _, entry := range seedTimeline() {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO timeline_entries (period, title, description) VALUES (?, ?, ?)`,
			entry.Period,
			entry.Title,
			entry.Description,
		); err != nil {
			return err
		}
	}
	return nil
}

func seedVideoRows(ctx context.Context, tx *sql.Tx) error {
	if hasRows, err := hasAnyRows(ctx, tx, "SELECT COUNT(1) FROM videos"); err != nil {
		return err
	} else if hasRows {
		return nil
	}

	for _, video := range seedVideos() {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO videos (title, description, url, thumbnail_url, published_at) VALUES (?, ?, ?, ?, ?)`,
			video.Title,
			video.Description,
			video.URL,
			video.ThumbnailURL,
			video.PublishedAt.Format(timeLayout),
		); err != nil {
			return err
		}
	}
	return nil
}

func hasAnyRows(ctx context.Context, tx *sql.Tx, query string) (bool, error) {
	var count int
	if err := tx.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

const timeLayout = time.RFC3339
