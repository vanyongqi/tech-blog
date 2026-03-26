package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"personal/blog/backend/model"
)

type SQLiteRepository struct {
	db *sql.DB
}

var ErrDuplicatePostSlug = errors.New("duplicate post slug")

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) GetSiteProfile(ctx context.Context) (model.SiteProfile, error) {
	var site model.SiteProfile
	var techStackJSON string
	var statsJSON string
	var socialLinksJSON string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT name, headline, intro, location, domain, email, motto, tech_stack_json, stats_json, social_links_json
		FROM site_profile
		WHERE id = 1`,
	).Scan(
		&site.Name,
		&site.Headline,
		&site.Intro,
		&site.Location,
		&site.Domain,
		&site.Email,
		&site.Motto,
		&techStackJSON,
		&statsJSON,
		&socialLinksJSON,
	)
	if err != nil {
		return model.SiteProfile{}, err
	}

	if err := json.Unmarshal([]byte(techStackJSON), &site.TechStack); err != nil {
		return model.SiteProfile{}, err
	}
	if err := json.Unmarshal([]byte(statsJSON), &site.Stats); err != nil {
		return model.SiteProfile{}, err
	}
	if err := json.Unmarshal([]byte(socialLinksJSON), &site.SocialLinks); err != nil {
		return model.SiteProfile{}, err
	}

	return site, nil
}

func (r *SQLiteRepository) ListPosts(ctx context.Context) ([]model.Post, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			p.id,
			p.slug,
			p.title,
			p.summary,
			p.category,
			p.read_time,
			p.hero_note,
			p.cover_label,
			p.tags_json,
			p.featured,
			p.published_at,
			p.blocks_json,
			(SELECT COUNT(1) FROM post_likes pl WHERE pl.post_id = p.id) AS like_count,
			(SELECT COUNT(1) FROM comments c WHERE c.post_id = p.id AND c.status = 'approved') AS comment_count
		FROM posts p
		ORDER BY p.published_at DESC, p.id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]model.Post, 0)
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func (r *SQLiteRepository) GetPostBySlug(ctx context.Context, slug, visitorID string) (model.Post, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT
			p.id,
			p.slug,
			p.title,
			p.summary,
			p.category,
			p.read_time,
			p.hero_note,
			p.cover_label,
			p.tags_json,
			p.featured,
			p.published_at,
			p.blocks_json,
			(SELECT COUNT(1) FROM post_likes pl WHERE pl.post_id = p.id) AS like_count,
			(SELECT COUNT(1) FROM comments c WHERE c.post_id = p.id AND c.status = 'approved') AS comment_count,
			EXISTS(SELECT 1 FROM post_likes pl WHERE pl.post_id = p.id AND pl.visitor_id = ?)
		FROM posts p
		WHERE p.slug = ?`,
		visitorID,
		slug,
	)

	post, err := scanPostWithLikeState(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Post{}, ErrPostNotFound
		}
		return model.Post{}, err
	}

	comments, err := r.listCommentsByPostID(ctx, post.ID)
	if err != nil {
		return model.Post{}, err
	}
	post.Comments = comments
	return post, nil
}

func (r *SQLiteRepository) CreatePost(ctx context.Context, post model.Post) (model.Post, error) {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO posts
			(slug, title, summary, category, read_time, hero_note, cover_label, tags_json, featured, published_at, blocks_json)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		post.Slug,
		post.Title,
		post.Summary,
		post.Category,
		post.ReadTime,
		post.HeroNote,
		post.CoverLabel,
		mustJSON(post.Tags),
		boolToInt(post.Featured),
		post.PublishedAt.Format(timeLayout),
		mustJSON(post.Blocks),
	)
	if err != nil {
		if isDuplicateSlugError(err) {
			return model.Post{}, ErrDuplicatePostSlug
		}
		return model.Post{}, err
	}
	return r.GetPostBySlug(ctx, post.Slug, "")
}

func (r *SQLiteRepository) UpdatePost(ctx context.Context, currentSlug string, post model.Post) (model.Post, error) {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE posts
		SET slug = ?, title = ?, summary = ?, category = ?, read_time = ?, hero_note = ?, cover_label = ?, tags_json = ?, featured = ?, published_at = ?, blocks_json = ?
		WHERE slug = ?`,
		post.Slug,
		post.Title,
		post.Summary,
		post.Category,
		post.ReadTime,
		post.HeroNote,
		post.CoverLabel,
		mustJSON(post.Tags),
		boolToInt(post.Featured),
		post.PublishedAt.Format(timeLayout),
		mustJSON(post.Blocks),
		currentSlug,
	)
	if err != nil {
		if isDuplicateSlugError(err) {
			return model.Post{}, ErrDuplicatePostSlug
		}
		return model.Post{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Post{}, err
	}
	if rowsAffected == 0 {
		return model.Post{}, ErrPostNotFound
	}
	return r.GetPostBySlug(ctx, post.Slug, "")
}

func (r *SQLiteRepository) DeletePost(ctx context.Context, slug string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	postID, err := r.lookupPostID(ctx, tx, slug)
	if err != nil {
		return err
	}

	for _, query := range []string{
		`DELETE FROM comments WHERE post_id = ?`,
		`DELETE FROM post_likes WHERE post_id = ?`,
		`DELETE FROM posts WHERE id = ?`,
	} {
		if _, err := tx.ExecContext(ctx, query, postID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteRepository) ListProjects(ctx context.Context) ([]model.Project, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, summary, status, link, image_url, accent, tech_stack_json FROM projects ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]model.Project, 0)
	for rows.Next() {
		project, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, rows.Err()
}

func (r *SQLiteRepository) GetProjectByID(ctx context.Context, id int64) (model.Project, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, summary, status, link, image_url, accent, tech_stack_json FROM projects WHERE id = ?`, id)
	project, err := scanProject(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Project{}, ErrProjectNotFound
		}
		return model.Project{}, err
	}
	return project, nil
}

func (r *SQLiteRepository) CreateProject(ctx context.Context, project model.Project) (model.Project, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO projects (name, summary, status, link, image_url, accent, tech_stack_json) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		project.Name,
		project.Summary,
		project.Status,
		project.Link,
		project.ImageURL,
		project.Accent,
		mustJSON(project.TechStack),
	)
	if err != nil {
		return model.Project{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Project{}, err
	}
	return r.GetProjectByID(ctx, id)
}

func (r *SQLiteRepository) UpdateProject(ctx context.Context, id int64, project model.Project) (model.Project, error) {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE projects SET name = ?, summary = ?, status = ?, link = ?, image_url = ?, accent = ?, tech_stack_json = ? WHERE id = ?`,
		project.Name,
		project.Summary,
		project.Status,
		project.Link,
		project.ImageURL,
		project.Accent,
		mustJSON(project.TechStack),
		id,
	)
	if err != nil {
		return model.Project{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Project{}, err
	}
	if rowsAffected == 0 {
		return model.Project{}, ErrProjectNotFound
	}
	return r.GetProjectByID(ctx, id)
}

func (r *SQLiteRepository) DeleteProject(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM projects WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrProjectNotFound
	}
	return nil
}

func (r *SQLiteRepository) ListVideos(ctx context.Context) ([]model.Video, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, title, description, url, thumbnail_url, published_at FROM videos ORDER BY published_at DESC, id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := make([]model.Video, 0)
	for rows.Next() {
		video, err := scanVideo(rows)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, rows.Err()
}

func (r *SQLiteRepository) GetVideoByID(ctx context.Context, id int64) (model.Video, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, title, description, url, thumbnail_url, published_at FROM videos WHERE id = ?`, id)
	video, err := scanVideo(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Video{}, ErrVideoNotFound
		}
		return model.Video{}, err
	}
	return video, nil
}

func (r *SQLiteRepository) CreateVideo(ctx context.Context, video model.Video) (model.Video, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO videos (title, description, url, thumbnail_url, published_at) VALUES (?, ?, ?, ?, ?)`,
		video.Title,
		video.Description,
		video.URL,
		video.ThumbnailURL,
		video.PublishedAt.Format(timeLayout),
	)
	if err != nil {
		return model.Video{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Video{}, err
	}
	return r.GetVideoByID(ctx, id)
}

func (r *SQLiteRepository) UpdateVideo(ctx context.Context, id int64, video model.Video) (model.Video, error) {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE videos SET title = ?, description = ?, url = ?, thumbnail_url = ?, published_at = ? WHERE id = ?`,
		video.Title,
		video.Description,
		video.URL,
		video.ThumbnailURL,
		video.PublishedAt.Format(timeLayout),
		id,
	)
	if err != nil {
		return model.Video{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Video{}, err
	}
	if rowsAffected == 0 {
		return model.Video{}, ErrVideoNotFound
	}
	return r.GetVideoByID(ctx, id)
}

func (r *SQLiteRepository) DeleteVideo(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM videos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrVideoNotFound
	}
	return nil
}

func (r *SQLiteRepository) ListTimeline(ctx context.Context) ([]model.TimelineEntry, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT period, title, description FROM timeline_entries ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]model.TimelineEntry, 0)
	for rows.Next() {
		var entry model.TimelineEntry
		if err := rows.Scan(&entry.Period, &entry.Title, &entry.Description); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (r *SQLiteRepository) CreateComment(ctx context.Context, input model.CreateCommentInput) (model.Comment, error) {
	postID, err := r.lookupPostID(ctx, r.db, input.Slug)
	if err != nil {
		return model.Comment{}, err
	}

	comment := model.Comment{
		AuthorName: input.AuthorName,
		Content:    input.Content,
		CreatedAt:  time.Now().UTC(),
	}

	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO comments (post_id, author_name, content, status, created_at)
		VALUES (?, ?, ?, 'approved', ?)`,
		postID,
		comment.AuthorName,
		comment.Content,
		comment.CreatedAt.Format(timeLayout),
	)
	if err != nil {
		return model.Comment{}, err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return model.Comment{}, err
	}
	comment.ID = commentID
	return comment, nil
}

func (r *SQLiteRepository) ToggleLike(ctx context.Context, input model.ToggleLikeInput) (model.LikeState, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.LikeState{}, err
	}
	defer tx.Rollback()

	postID, err := r.lookupPostID(ctx, tx, input.Slug)
	if err != nil {
		return model.LikeState{}, err
	}

	result, err := tx.ExecContext(
		ctx,
		`DELETE FROM post_likes WHERE post_id = ? AND visitor_id = ?`,
		postID,
		input.VisitorID,
	)
	if err != nil {
		return model.LikeState{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.LikeState{}, err
	}

	liked := false
	if rowsAffected == 0 {
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO post_likes (post_id, visitor_id, created_at) VALUES (?, ?, ?)`,
			postID,
			input.VisitorID,
			time.Now().UTC().Format(timeLayout),
		); err != nil {
			return model.LikeState{}, err
		}
		liked = true
	}

	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM post_likes WHERE post_id = ?`, postID).Scan(&count); err != nil {
		return model.LikeState{}, err
	}

	if err := tx.Commit(); err != nil {
		return model.LikeState{}, err
	}

	return model.LikeState{
		LikeCount: count,
		Liked:     liked,
	}, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPost(scanner rowScanner) (model.Post, error) {
	var post model.Post
	var tagsJSON string
	var blocksJSON string
	var featuredInt int
	var publishedAt string

	if err := scanner.Scan(
		&post.ID,
		&post.Slug,
		&post.Title,
		&post.Summary,
		&post.Category,
		&post.ReadTime,
		&post.HeroNote,
		&post.CoverLabel,
		&tagsJSON,
		&featuredInt,
		&publishedAt,
		&blocksJSON,
		&post.LikeCount,
		&post.CommentCount,
	); err != nil {
		return model.Post{}, err
	}

	if err := hydratePost(&post, tagsJSON, blocksJSON, featuredInt, publishedAt); err != nil {
		return model.Post{}, err
	}
	return post, nil
}

func scanPostWithLikeState(scanner rowScanner) (model.Post, error) {
	var post model.Post
	var tagsJSON string
	var blocksJSON string
	var featuredInt int
	var publishedAt string
	var likedByVisitorInt int

	if err := scanner.Scan(
		&post.ID,
		&post.Slug,
		&post.Title,
		&post.Summary,
		&post.Category,
		&post.ReadTime,
		&post.HeroNote,
		&post.CoverLabel,
		&tagsJSON,
		&featuredInt,
		&publishedAt,
		&blocksJSON,
		&post.LikeCount,
		&post.CommentCount,
		&likedByVisitorInt,
	); err != nil {
		return model.Post{}, err
	}

	if err := hydratePost(&post, tagsJSON, blocksJSON, featuredInt, publishedAt); err != nil {
		return model.Post{}, err
	}
	post.LikedByVisitor = likedByVisitorInt == 1
	return post, nil
}

func hydratePost(post *model.Post, tagsJSON, blocksJSON string, featuredInt int, publishedAt string) error {
	if err := json.Unmarshal([]byte(tagsJSON), &post.Tags); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(blocksJSON), &post.Blocks); err != nil {
		return err
	}

	parsedTime, err := time.Parse(timeLayout, publishedAt)
	if err != nil {
		return err
	}
	post.PublishedAt = parsedTime
	post.Featured = featuredInt == 1
	return nil
}

func (r *SQLiteRepository) listCommentsByPostID(ctx context.Context, postID int64) ([]model.Comment, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, author_name, content, created_at
		FROM comments
		WHERE post_id = ? AND status = 'approved'
		ORDER BY created_at DESC`,
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]model.Comment, 0)
	for rows.Next() {
		var comment model.Comment
		var createdAt string
		if err := rows.Scan(&comment.ID, &comment.AuthorName, &comment.Content, &createdAt); err != nil {
			return nil, err
		}
		parsedTime, err := time.Parse(timeLayout, createdAt)
		if err != nil {
			return nil, err
		}
		comment.CreatedAt = parsedTime
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func scanVideo(scanner rowScanner) (model.Video, error) {
	var video model.Video
	var publishedAt string
	if err := scanner.Scan(&video.ID, &video.Title, &video.Description, &video.URL, &video.ThumbnailURL, &publishedAt); err != nil {
		return model.Video{}, err
	}
	parsedTime, err := time.Parse(timeLayout, publishedAt)
	if err != nil {
		return model.Video{}, err
	}
	video.PublishedAt = parsedTime
	return video, nil
}

func scanProject(scanner rowScanner) (model.Project, error) {
	var project model.Project
	var techStackJSON string
	if err := scanner.Scan(
		&project.ID,
		&project.Name,
		&project.Summary,
		&project.Status,
		&project.Link,
		&project.ImageURL,
		&project.Accent,
		&techStackJSON,
	); err != nil {
		return model.Project{}, err
	}
	if err := json.Unmarshal([]byte(techStackJSON), &project.TechStack); err != nil {
		return model.Project{}, err
	}
	return project, nil
}

type postLookupQuery interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (r *SQLiteRepository) lookupPostID(ctx context.Context, queryer postLookupQuery, slug string) (int64, error) {
	var postID int64
	err := queryer.QueryRowContext(ctx, `SELECT id FROM posts WHERE slug = ?`, slug).Scan(&postID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrPostNotFound
		}
		return 0, err
	}
	return postID, nil
}

func isDuplicateSlugError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint failed: posts.slug")
}
