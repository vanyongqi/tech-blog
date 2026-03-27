package model

import "time"

type SiteStat struct {
	Label string
	Value string
}

type SocialLink struct {
	Label string
	URL   string
}

type SiteProfile struct {
	Name        string
	Headline    string
	Intro       string
	Location    string
	Domain      string
	Email       string
	Motto       string
	TechStack   []string
	Stats       []SiteStat
	SocialLinks []SocialLink
}

type Post struct {
	ID             int64
	Slug           string
	Title          string
	Summary        string
	Category       string
	ReadTime       string
	CoverLabel     string
	ContentMarkdown string
	Tags           []string
	Featured       bool
	PublishedAt    time.Time
	LikeCount      int
	CommentCount   int
	LikedByVisitor bool
	Comments       []Comment
}

type Asset struct {
	ID        int64
	Filename  string
	MimeType  string
	Data      []byte
	CreatedAt time.Time
}

type Project struct {
	ID        int64
	Name      string
	Summary   string
	Status    string
	Link      string
	ImageURL  string
	Accent    string
	TechStack []string
}

type Video struct {
	ID          int64
	Title       string
	Description string
	URL         string
	ThumbnailURL string
	PublishedAt time.Time
}

type TimelineEntry struct {
	Period      string
	Title       string
	Description string
}

type Comment struct {
	ID         int64
	AuthorName string
	Content    string
	CreatedAt  time.Time
}

type LikeState struct {
	LikeCount int
	Liked     bool
}

type ListPostsInput struct {
	FeaturedOnly bool
	Tag          string
	Limit        int
}

type GetPostInput struct {
	Slug      string
	VisitorID string
}

type CreateCommentInput struct {
	Slug       string
	AuthorName string
	Content    string
}

type ToggleLikeInput struct {
	Slug      string
	VisitorID string
}
