package service

import (
	"context"
	"errors"
	"slices"
	"strings"

	"personal/blog/backend/dao"
	"personal/blog/backend/model"
)

var ErrPostNotFound = errors.New("post not found")
var ErrCommentContentRequired = errors.New("comment content is required")
var ErrCommentContentTooLong = errors.New("comment content must be 500 characters or fewer")
var ErrVisitorIdentityRequired = errors.New("visitor identity is required")
var ErrAssetNotFound = errors.New("asset not found")

type BlogService struct {
	repository dao.BlogRepository
}

func NewBlogService(repository dao.BlogRepository) *BlogService {
	return &BlogService{repository: repository}
}

func (s *BlogService) GetSiteProfile(ctx context.Context) (model.SiteProfile, error) {
	return s.repository.GetSiteProfile(ctx)
}

func (s *BlogService) ListPosts(ctx context.Context, input model.ListPostsInput) ([]model.Post, error) {
	posts, err := s.repository.ListPosts(ctx)
	if err != nil {
		return nil, err
	}
	filtered := make([]model.Post, 0, len(posts))
	tag := strings.TrimSpace(strings.ToLower(input.Tag))

	for _, post := range posts {
		if input.FeaturedOnly && !post.Featured {
			continue
		}
		if tag != "" && !containsTag(post.Tags, tag) {
			continue
		}
		filtered = append(filtered, post)
	}

	slices.SortFunc(filtered, func(a, b model.Post) int {
		switch {
		case a.PublishedAt.After(b.PublishedAt):
			return -1
		case a.PublishedAt.Before(b.PublishedAt):
			return 1
		default:
			return strings.Compare(a.Title, b.Title)
		}
	})

	if input.Limit > 0 && input.Limit < len(filtered) {
		return filtered[:input.Limit], nil
	}

	return filtered, nil
}

func (s *BlogService) GetPost(ctx context.Context, input model.GetPostInput) (model.Post, error) {
	post, err := s.repository.GetPostBySlug(ctx, input.Slug, input.VisitorID)
	if err != nil {
		if errors.Is(err, dao.ErrPostNotFound) {
			return model.Post{}, ErrPostNotFound
		}
		return model.Post{}, err
	}
	return post, nil
}

func (s *BlogService) ListProjects(ctx context.Context) ([]model.Project, error) {
	return s.repository.ListProjects(ctx)
}

func (s *BlogService) ListVideos(ctx context.Context) ([]model.Video, error) {
	return s.repository.ListVideos(ctx)
}

func (s *BlogService) ListTimeline(ctx context.Context) ([]model.TimelineEntry, error) {
	return s.repository.ListTimeline(ctx)
}

func (s *BlogService) GetAsset(ctx context.Context, id int64) (model.Asset, error) {
	asset, err := s.repository.GetAssetByID(ctx, id)
	if err != nil {
		if errors.Is(err, dao.ErrAssetNotFound) {
			return model.Asset{}, ErrAssetNotFound
		}
		return model.Asset{}, err
	}
	return asset, nil
}

func (s *BlogService) CreateComment(ctx context.Context, input model.CreateCommentInput) (model.Comment, error) {
	input.AuthorName = strings.TrimSpace(input.AuthorName)
	input.Content = strings.TrimSpace(input.Content)

	if input.AuthorName == "" {
		return model.Comment{}, ErrVisitorIdentityRequired
	}
	if input.Content == "" {
		return model.Comment{}, ErrCommentContentRequired
	}
	if len([]rune(input.Content)) > 500 {
		return model.Comment{}, ErrCommentContentTooLong
	}

	comment, err := s.repository.CreateComment(ctx, input)
	if err != nil {
		if errors.Is(err, dao.ErrPostNotFound) {
			return model.Comment{}, ErrPostNotFound
		}
		return model.Comment{}, err
	}
	return comment, nil
}

func (s *BlogService) ToggleLike(ctx context.Context, input model.ToggleLikeInput) (model.LikeState, error) {
	input.VisitorID = strings.TrimSpace(input.VisitorID)
	if input.VisitorID == "" {
		return model.LikeState{}, ErrVisitorIdentityRequired
	}

	state, err := s.repository.ToggleLike(ctx, input)
	if err != nil {
		if errors.Is(err, dao.ErrPostNotFound) {
			return model.LikeState{}, ErrPostNotFound
		}
		return model.LikeState{}, err
	}
	return state, nil
}

func containsTag(tags []string, target string) bool {
	for _, tag := range tags {
		if strings.EqualFold(tag, target) {
			return true
		}
	}
	return false
}
