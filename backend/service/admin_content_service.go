package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"personal/blog/backend/dao"
	"personal/blog/backend/model"
)

var (
	ErrPostSlugRequired       = errors.New("post slug is required")
	ErrPostSlugInvalid        = errors.New("post slug must contain lowercase letters, numbers or hyphens")
	ErrPostTitleRequired      = errors.New("post title is required")
	ErrPostSummaryRequired    = errors.New("post summary is required")
	ErrPostCategoryRequired   = errors.New("post category is required")
	ErrPostReadTimeRequired   = errors.New("post read time is required")
	ErrPostPublishedAtInvalid = errors.New("post published date is invalid")
	ErrPostBlocksRequired     = errors.New("post content blocks are required")
	ErrDuplicatePostSlug      = errors.New("post slug already exists")
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type AdminContentService struct {
	repository dao.BlogRepository
}

func NewAdminContentService(repository dao.BlogRepository) *AdminContentService {
	return &AdminContentService{repository: repository}
}

func (s *AdminContentService) ListPosts(ctx context.Context) ([]model.Post, error) {
	return s.repository.ListPosts(ctx)
}

func (s *AdminContentService) GetPost(ctx context.Context, input model.GetPostInput) (model.Post, error) {
	post, err := s.repository.GetPostBySlug(ctx, input.Slug, input.VisitorID)
	if err != nil {
		if errors.Is(err, dao.ErrPostNotFound) {
			return model.Post{}, ErrPostNotFound
		}
		return model.Post{}, err
	}
	return post, nil
}

func (s *AdminContentService) CreatePost(ctx context.Context, input model.CreatePostInput) (model.Post, error) {
	post, err := normalizeAdminPost(input.Post)
	if err != nil {
		return model.Post{}, err
	}

	created, err := s.repository.CreatePost(ctx, post)
	if err != nil {
		if errors.Is(err, dao.ErrDuplicatePostSlug) {
			return model.Post{}, ErrDuplicatePostSlug
		}
		return model.Post{}, err
	}
	return created, nil
}

func (s *AdminContentService) UpdatePost(ctx context.Context, input model.UpdatePostInput) (model.Post, error) {
	post, err := normalizeAdminPost(input.Post)
	if err != nil {
		return model.Post{}, err
	}

	updated, err := s.repository.UpdatePost(ctx, strings.TrimSpace(input.CurrentSlug), post)
	if err != nil {
		switch {
		case errors.Is(err, dao.ErrPostNotFound):
			return model.Post{}, ErrPostNotFound
		case errors.Is(err, dao.ErrDuplicatePostSlug):
			return model.Post{}, ErrDuplicatePostSlug
		default:
			return model.Post{}, err
		}
	}
	return updated, nil
}

func (s *AdminContentService) DeletePost(ctx context.Context, input model.DeletePostInput) error {
	err := s.repository.DeletePost(ctx, strings.TrimSpace(input.Slug))
	if err != nil {
		if errors.Is(err, dao.ErrPostNotFound) {
			return ErrPostNotFound
		}
		return err
	}
	return nil
}

func normalizeAdminPost(post model.Post) (model.Post, error) {
	post.Slug = strings.TrimSpace(post.Slug)
	post.Title = strings.TrimSpace(post.Title)
	post.Summary = strings.TrimSpace(post.Summary)
	post.Category = strings.TrimSpace(post.Category)
	post.ReadTime = strings.TrimSpace(post.ReadTime)
	post.HeroNote = strings.TrimSpace(post.HeroNote)
	post.CoverLabel = strings.TrimSpace(post.CoverLabel)

	if post.Slug == "" {
		return model.Post{}, ErrPostSlugRequired
	}
	if !slugPattern.MatchString(post.Slug) {
		return model.Post{}, ErrPostSlugInvalid
	}
	if post.Title == "" {
		return model.Post{}, ErrPostTitleRequired
	}
	if post.Summary == "" {
		return model.Post{}, ErrPostSummaryRequired
	}
	if post.Category == "" {
		return model.Post{}, ErrPostCategoryRequired
	}
	if post.ReadTime == "" {
		return model.Post{}, ErrPostReadTimeRequired
	}
	if post.PublishedAt.IsZero() {
		return model.Post{}, ErrPostPublishedAtInvalid
	}

	post.Tags = normalizeTags(post.Tags)
	post.Blocks = normalizeBlocks(post.Blocks)
	if len(post.Blocks) == 0 {
		return model.Post{}, ErrPostBlocksRequired
	}

	return post, nil
}

func normalizeTags(tags []string) []string {
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		normalized := strings.TrimSpace(tag)
		if normalized == "" {
			continue
		}
		result = append(result, normalized)
	}
	return result
}

func normalizeBlocks(blocks []model.ContentBlock) []model.ContentBlock {
	result := make([]model.ContentBlock, 0, len(blocks))
	for _, block := range blocks {
		normalized := model.ContentBlock{
			Kind:  strings.TrimSpace(block.Kind),
			Title: strings.TrimSpace(block.Title),
			Text:  strings.TrimSpace(block.Text),
			URL:   strings.TrimSpace(block.URL),
		}

		if normalized.Kind == "" {
			continue
		}

		if len(block.Items) > 0 {
			items := make([]string, 0, len(block.Items))
			for _, item := range block.Items {
				trimmed := strings.TrimSpace(item)
				if trimmed == "" {
					continue
				}
				items = append(items, trimmed)
			}
			normalized.Items = items
		}

		if normalized.Kind == "list" && len(normalized.Items) == 0 {
			continue
		}
		if normalized.Kind == "video" && normalized.URL == "" {
			continue
		}
		if normalized.Kind != "list" && normalized.Kind != "video" && normalized.Text == "" && normalized.Title == "" {
			continue
		}

		result = append(result, normalized)
	}
	return result
}
