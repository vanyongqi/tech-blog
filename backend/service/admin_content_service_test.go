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
			Title:       "后台新建文章",
			Summary:     "通过管理后台创建的新文章。",
			Category:    "Product",
			ReadTime:    "4 分钟",
			CoverLabel:  "后台内容",
			ContentMarkdown: "# 新文章\n\n这是一段新内容。",
			Tags:        []string{"Admin", "React"},
			PublishedAt: time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.Slug == "" {
		t.Fatal("expected generated slug")
	}
}

func TestAdminContentServiceValidatesPost(t *testing.T) {
	svc := NewAdminContentService(newStubRepository())

	_, err := svc.CreatePost(context.Background(), model.CreatePostInput{
		Post: model.Post{
			Slug:            "Bad Slug",
			Title:           "",
			Summary:         "",
			Category:        "",
			ReadTime:        "",
			ContentMarkdown: "",
			PublishedAt:     time.Time{},
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
			CoverLabel:  "系统设计",
			ContentMarkdown: "## 更新后的正文\n\n这里是新的 Markdown 内容。",
			Tags:        []string{"Golang"},
			Featured:    true,
			PublishedAt: time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC),
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

func TestAdminContentServiceCreateAsset(t *testing.T) {
	svc := NewAdminContentService(newStubRepository())

	asset, err := svc.CreateAsset(context.Background(), model.CreateAssetInput{
		Asset: model.Asset{
			Filename: "diagram.png",
			MimeType: "image/png",
			Data:     []byte("fake-image-data"),
		},
	})
	if err != nil {
		t.Fatalf("unexpected asset error: %v", err)
	}
	if asset.ID == 0 {
		t.Fatal("expected asset id to be generated")
	}
}
