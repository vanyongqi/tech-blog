package service

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"personal/blog/backend/dao"
	"personal/blog/backend/model"
)

var (
	ErrPostSlugInvalid        = errors.New("post slug must contain lowercase letters, numbers or hyphens")
	ErrPostTitleRequired      = errors.New("post title is required")
	ErrPostSummaryRequired    = errors.New("post summary is required")
	ErrPostCategoryRequired   = errors.New("post category is required")
	ErrPostReadTimeRequired   = errors.New("post read time is required")
	ErrPostPublishedAtInvalid = errors.New("post published date is invalid")
	ErrPostContentMarkdownRequired = errors.New("post markdown content is required")
	ErrDuplicatePostSlug      = errors.New("post slug already exists")
	ErrAssetFileRequired      = errors.New("image file is required")
	ErrAssetMimeTypeInvalid   = errors.New("only image uploads are supported")
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
	post, err := normalizeAdminPost(input.Post, "")
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
	post, err := normalizeAdminPost(input.Post, strings.TrimSpace(input.CurrentSlug))
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

func (s *AdminContentService) CreateAsset(ctx context.Context, input model.CreateAssetInput) (model.Asset, error) {
	asset, err := normalizeAdminAsset(input.Asset)
	if err != nil {
		return model.Asset{}, err
	}
	return s.repository.CreateAsset(ctx, asset)
}

func normalizeAdminPost(post model.Post, fallbackSlug string) (model.Post, error) {
	post.Slug = strings.TrimSpace(post.Slug)
	post.Title = strings.TrimSpace(post.Title)
	post.Summary = strings.TrimSpace(post.Summary)
	post.Category = strings.TrimSpace(post.Category)
	post.ReadTime = strings.TrimSpace(post.ReadTime)
	post.CoverLabel = strings.TrimSpace(post.CoverLabel)
	post.ContentMarkdown = strings.TrimSpace(post.ContentMarkdown)

	if post.Slug == "" {
		post.Slug = strings.TrimSpace(fallbackSlug)
	}
	if post.Slug == "" {
		post.Slug = generatePostSlug()
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
	if post.ContentMarkdown == "" {
		return model.Post{}, ErrPostContentMarkdownRequired
	}

	post.Tags = normalizeTags(post.Tags)

	return post, nil
}

func normalizeAdminAsset(asset model.Asset) (model.Asset, error) {
	asset.Filename = strings.TrimSpace(asset.Filename)
	asset.MimeType = strings.TrimSpace(asset.MimeType)

	if len(asset.Data) == 0 {
		return model.Asset{}, ErrAssetFileRequired
	}
	if asset.MimeType == "" {
		asset.MimeType = http.DetectContentType(asset.Data)
	}
	if !strings.HasPrefix(asset.MimeType, "image/") {
		return model.Asset{}, ErrAssetMimeTypeInvalid
	}
	if asset.Filename == "" {
		asset.Filename = "image"
	}
	return asset, nil
}

func generatePostSlug() string {
	return "post-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
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
