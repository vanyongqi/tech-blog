package dao

import (
	"context"
	"errors"

	"personal/blog/backend/model"
)

var ErrPostNotFound = errors.New("post not found")
var ErrProjectNotFound = errors.New("project not found")
var ErrVideoNotFound = errors.New("video not found")

type BlogRepository interface {
	GetSiteProfile(ctx context.Context) (model.SiteProfile, error)
	ListPosts(ctx context.Context) ([]model.Post, error)
	GetPostBySlug(ctx context.Context, slug, visitorID string) (model.Post, error)
	CreatePost(ctx context.Context, post model.Post) (model.Post, error)
	UpdatePost(ctx context.Context, currentSlug string, post model.Post) (model.Post, error)
	DeletePost(ctx context.Context, slug string) error
	ListProjects(ctx context.Context) ([]model.Project, error)
	GetProjectByID(ctx context.Context, id int64) (model.Project, error)
	CreateProject(ctx context.Context, project model.Project) (model.Project, error)
	UpdateProject(ctx context.Context, id int64, project model.Project) (model.Project, error)
	DeleteProject(ctx context.Context, id int64) error
	ListVideos(ctx context.Context) ([]model.Video, error)
	GetVideoByID(ctx context.Context, id int64) (model.Video, error)
	CreateVideo(ctx context.Context, video model.Video) (model.Video, error)
	UpdateVideo(ctx context.Context, id int64, video model.Video) (model.Video, error)
	DeleteVideo(ctx context.Context, id int64) error
	ListTimeline(ctx context.Context) ([]model.TimelineEntry, error)
	CreateComment(ctx context.Context, input model.CreateCommentInput) (model.Comment, error)
	ToggleLike(ctx context.Context, input model.ToggleLikeInput) (model.LikeState, error)
}
