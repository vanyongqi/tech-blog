package service

import (
	"context"
	"testing"
	"time"

	"personal/blog/backend/model"
)

func TestAdminContentServiceCreatePost(t *testing.T) {
	svc := NewAdminContentService(newStubRepository())

	post, err := svc.CreatePost(context.Background(), model.CreatePostInput{
		Post: model.Post{
			Slug:        "new-admin-post",
			Title:       "后台新建文章",
			Summary:     "通过管理后台创建的新文章。",
			Category:    "Product",
			ReadTime:    "4 分钟",
			HeroNote:    "让内容维护进入日常流程。",
			CoverLabel:  "后台内容",
			Tags:        []string{"Admin", "React"},
			PublishedAt: time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC),
			Blocks: []model.ContentBlock{
				{Kind: "paragraph", Text: "这是一段新内容。"},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.Slug != "new-admin-post" {
		t.Fatalf("expected slug new-admin-post, got %q", post.Slug)
	}
}

func TestAdminContentServiceValidatesPost(t *testing.T) {
	svc := NewAdminContentService(newStubRepository())

	_, err := svc.CreatePost(context.Background(), model.CreatePostInput{
		Post: model.Post{
			Slug:        "Bad Slug",
			Title:       "",
			Summary:     "",
			Category:    "",
			ReadTime:    "",
			PublishedAt: time.Time{},
		},
	})
	if err != ErrPostSlugInvalid {
		t.Fatalf("expected %v, got %v", ErrPostSlugInvalid, err)
	}
}

func TestAdminContentServiceUpdateAndDeletePost(t *testing.T) {
	svc := NewAdminContentService(newStubRepository())
	ctx := context.Background()

	updated, err := svc.UpdatePost(ctx, model.UpdatePostInput{
		CurrentSlug: "designing-calm-backends",
		Post: model.Post{
			Slug:        "designing-calm-backends-v2",
			Title:       "如何把后端系统做得更稳定",
			Summary:     "更新后的摘要。",
			Category:    "Engineering",
			ReadTime:    "9 分钟",
			HeroNote:    "更新后的引导语。",
			CoverLabel:  "系统设计",
			Tags:        []string{"Golang"},
			Featured:    true,
			PublishedAt: time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC),
			Blocks: []model.ContentBlock{
				{Kind: "paragraph", Text: "更新后的正文。"},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}
	if updated.Slug != "designing-calm-backends-v2" {
		t.Fatalf("expected updated slug, got %q", updated.Slug)
	}

	if err := svc.DeletePost(ctx, model.DeletePostInput{Slug: updated.Slug}); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
}
