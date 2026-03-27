package model

import "time"

type AdminLoginInput struct {
	Username string
	Password string
}

type AdminSession struct {
	Username  string
	ExpiresAt time.Time
}

type CreatePostInput struct {
	Post Post
}

type UpdatePostInput struct {
	CurrentSlug string
	Post        Post
}

type DeletePostInput struct {
	Slug string
}

type CreateAssetInput struct {
	Asset Asset
}

type CreateProjectInput struct {
	Project Project
}

type UpdateProjectInput struct {
	ID      int64
	Project Project
}

type DeleteProjectInput struct {
	ID int64
}

type CreateVideoInput struct {
	Video Video
}

type UpdateVideoInput struct {
	ID    int64
	Video Video
}

type DeleteVideoInput struct {
	ID int64
}
