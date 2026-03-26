package service

import (
	"context"
	"testing"
	"time"

	"personal/blog/backend/dao"
	"personal/blog/backend/model"
)

func TestBlogServiceListPosts(t *testing.T) {
	svc := NewBlogService(newStubRepository())
	ctx := context.Background()

	testCases := []struct {
		name      string
		input     model.ListPostsInput
		wantCount int
		wantSlug  string
	}{
		{name: "returns sorted posts by published time", input: model.ListPostsInput{}, wantCount: 3, wantSlug: "designing-calm-backends"},
		{name: "filters featured posts", input: model.ListPostsInput{FeaturedOnly: true}, wantCount: 2, wantSlug: "designing-calm-backends"},
		{name: "filters by tag and limit", input: model.ListPostsInput{Tag: "react", Limit: 1}, wantCount: 1, wantSlug: "building-a-personal-operating-system"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := svc.ListPosts(ctx, tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tc.wantCount {
				t.Fatalf("expected %d posts, got %d", tc.wantCount, len(got))
			}
			if tc.wantCount > 0 && got[0].Slug != tc.wantSlug {
				t.Fatalf("expected first slug %q, got %q", tc.wantSlug, got[0].Slug)
			}
		})
	}
}

func TestBlogServiceGetPost(t *testing.T) {
	svc := NewBlogService(newStubRepository())
	ctx := context.Background()

	testCases := []struct {
		name     string
		input    model.GetPostInput
		wantErr  error
		wantSlug string
	}{
		{name: "returns post by slug", input: model.GetPostInput{Slug: "designing-calm-backends", VisitorID: "visitor-a"}, wantSlug: "designing-calm-backends"},
		{name: "returns not found when missing", input: model.GetPostInput{Slug: "missing-post", VisitorID: "visitor-a"}, wantErr: ErrPostNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := svc.GetPost(ctx, tc.input)
			if tc.wantErr != nil {
				if err != tc.wantErr {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Slug != tc.wantSlug {
				t.Fatalf("expected slug %q, got %q", tc.wantSlug, got.Slug)
			}
		})
	}
}

func TestBlogServiceCreateComment(t *testing.T) {
	svc := NewBlogService(newStubRepository())
	ctx := context.Background()

	testCases := []struct {
		name        string
		input       model.CreateCommentInput
		wantErr     error
		wantAuthor  string
		wantContent string
	}{
		{
			name: "creates comment with trimmed content",
			input: model.CreateCommentInput{
				Slug:       "designing-calm-backends",
				AuthorName: "访客-AAAA1111",
				Content:    "  写得很好，想看更多系统设计文章。  ",
			},
			wantAuthor:  "访客-AAAA1111",
			wantContent: "写得很好，想看更多系统设计文章。",
		},
		{
			name: "rejects empty content",
			input: model.CreateCommentInput{
				Slug:       "designing-calm-backends",
				AuthorName: "访客-AAAA1111",
				Content:    "   ",
			},
			wantErr: ErrCommentContentRequired,
		},
		{
			name: "rejects missing visitor identity",
			input: model.CreateCommentInput{
				Slug:    "designing-calm-backends",
				Content: "内容有效",
			},
			wantErr: ErrVisitorIdentityRequired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := svc.CreateComment(ctx, tc.input)
			if tc.wantErr != nil {
				if err != tc.wantErr {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.AuthorName != tc.wantAuthor {
				t.Fatalf("expected author %q, got %q", tc.wantAuthor, got.AuthorName)
			}
			if got.Content != tc.wantContent {
				t.Fatalf("expected content %q, got %q", tc.wantContent, got.Content)
			}
		})
	}
}

func TestBlogServiceToggleLike(t *testing.T) {
	svc := NewBlogService(newStubRepository())
	ctx := context.Background()

	first, err := svc.ToggleLike(ctx, model.ToggleLikeInput{
		Slug:      "designing-calm-backends",
		VisitorID: "visitor-like",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !first.Liked || first.LikeCount != 1 {
		t.Fatalf("expected first toggle to like with count 1, got %+v", first)
	}

	second, err := svc.ToggleLike(ctx, model.ToggleLikeInput{
		Slug:      "designing-calm-backends",
		VisitorID: "visitor-like",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if second.Liked || second.LikeCount != 0 {
		t.Fatalf("expected second toggle to unlike with count 0, got %+v", second)
	}
}

func TestBlogServiceReadOnlyLists(t *testing.T) {
	svc := NewBlogService(newStubRepository())
	ctx := context.Background()

	site, err := svc.GetSiteProfile(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting site profile: %v", err)
	}
	if site.Name != "Fitz" {
		t.Fatalf("expected site name Fitz, got %q", site.Name)
	}

	projects, err := svc.ListProjects(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing projects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	timeline, err := svc.ListTimeline(ctx)
	if err != nil {
		t.Fatalf("unexpected error listing timeline: %v", err)
	}
	if len(timeline) != 1 {
		t.Fatalf("expected 1 timeline entry, got %d", len(timeline))
	}
}

type stubRepository struct {
	posts    []model.Post
	comments map[string][]model.Comment
	likes    map[string]map[string]bool
	projects []model.Project
	videos   []model.Video
}

func newStubRepository() *stubRepository {
	return &stubRepository{
		posts: []model.Post{
			{ID: 1, Slug: "designing-calm-backends", Title: "如何把后端系统做得稳定、克制、可增长", Tags: []string{"Golang", "Architecture", "Observability"}, Featured: true, PublishedAt: time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)},
			{ID: 2, Slug: "building-a-personal-operating-system", Title: "把个人博客做成个人操作系统", Tags: []string{"React", "Brand", "Writing"}, Featured: true, PublishedAt: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
			{ID: 3, Slug: "react-layouts-with-editorial-rhythm", Title: "给 React 页面做出杂志感，不靠堆组件", Tags: []string{"React", "CSS", "Design"}, PublishedAt: time.Date(2026, 2, 27, 0, 0, 0, 0, time.UTC)},
		},
		comments: map[string][]model.Comment{},
		likes:    map[string]map[string]bool{},
		projects: []model.Project{{ID: 1, Name: "Project", Link: "https://example.com"}},
		videos:   []model.Video{{ID: 1, Title: "Video", URL: "https://example.com/video", PublishedAt: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)}},
	}
}

func (r *stubRepository) GetSiteProfile(context.Context) (model.SiteProfile, error) {
	return model.SiteProfile{Name: "Fitz"}, nil
}

func (r *stubRepository) ListPosts(context.Context) ([]model.Post, error) {
	posts := make([]model.Post, 0, len(r.posts))
	for _, post := range r.posts {
		post.CommentCount = len(r.comments[post.Slug])
		post.LikeCount = len(r.likes[post.Slug])
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *stubRepository) GetPostBySlug(_ context.Context, slug, visitorID string) (model.Post, error) {
	for _, post := range r.posts {
		if post.Slug != slug {
			continue
		}
		post.CommentCount = len(r.comments[post.Slug])
		post.LikeCount = len(r.likes[post.Slug])
		post.LikedByVisitor = r.likes[post.Slug][visitorID]
		post.Comments = append([]model.Comment(nil), r.comments[post.Slug]...)
		return post, nil
	}
	return model.Post{}, dao.ErrPostNotFound
}

func (r *stubRepository) CreatePost(_ context.Context, post model.Post) (model.Post, error) {
	if _, err := r.GetPostBySlug(context.Background(), post.Slug, ""); err == nil {
		return model.Post{}, dao.ErrDuplicatePostSlug
	}
	post.ID = int64(len(r.posts) + 1)
	r.posts = append(r.posts, post)
	return post, nil
}

func (r *stubRepository) UpdatePost(_ context.Context, currentSlug string, post model.Post) (model.Post, error) {
	targetIndex := -1
	for index, item := range r.posts {
		if item.Slug == currentSlug {
			targetIndex = index
			continue
		}
		if item.Slug == post.Slug {
			return model.Post{}, dao.ErrDuplicatePostSlug
		}
	}
	if targetIndex < 0 {
		return model.Post{}, dao.ErrPostNotFound
	}
	post.ID = r.posts[targetIndex].ID
	post.CommentCount = len(r.comments[currentSlug])
	post.LikeCount = len(r.likes[currentSlug])
	post.Comments = append([]model.Comment(nil), r.comments[currentSlug]...)
	r.posts[targetIndex] = post
	if currentSlug != post.Slug {
		r.comments[post.Slug] = r.comments[currentSlug]
		delete(r.comments, currentSlug)
		r.likes[post.Slug] = r.likes[currentSlug]
		delete(r.likes, currentSlug)
	}
	return post, nil
}

func (r *stubRepository) DeletePost(_ context.Context, slug string) error {
	for index, post := range r.posts {
		if post.Slug != slug {
			continue
		}
		r.posts = append(r.posts[:index], r.posts[index+1:]...)
		delete(r.comments, slug)
		delete(r.likes, slug)
		return nil
	}
	return dao.ErrPostNotFound
}

func (r *stubRepository) ListProjects(context.Context) ([]model.Project, error) {
	return append([]model.Project(nil), r.projects...), nil
}

func (r *stubRepository) GetProjectByID(_ context.Context, id int64) (model.Project, error) {
	for _, project := range r.projects {
		if project.ID == id {
			return project, nil
		}
	}
	return model.Project{}, dao.ErrProjectNotFound
}

func (r *stubRepository) CreateProject(_ context.Context, project model.Project) (model.Project, error) {
	project.ID = int64(len(r.projects) + 1)
	r.projects = append(r.projects, project)
	return project, nil
}

func (r *stubRepository) UpdateProject(_ context.Context, id int64, project model.Project) (model.Project, error) {
	for index, current := range r.projects {
		if current.ID == id {
			project.ID = id
			r.projects[index] = project
			return project, nil
		}
	}
	return model.Project{}, dao.ErrProjectNotFound
}

func (r *stubRepository) DeleteProject(_ context.Context, id int64) error {
	for index, project := range r.projects {
		if project.ID == id {
			r.projects = append(r.projects[:index], r.projects[index+1:]...)
			return nil
		}
	}
	return dao.ErrProjectNotFound
}

func (r *stubRepository) ListVideos(context.Context) ([]model.Video, error) {
	return append([]model.Video(nil), r.videos...), nil
}

func (r *stubRepository) GetVideoByID(_ context.Context, id int64) (model.Video, error) {
	for _, video := range r.videos {
		if video.ID == id {
			return video, nil
		}
	}
	return model.Video{}, dao.ErrVideoNotFound
}

func (r *stubRepository) CreateVideo(_ context.Context, video model.Video) (model.Video, error) {
	video.ID = int64(len(r.videos) + 1)
	r.videos = append(r.videos, video)
	return video, nil
}

func (r *stubRepository) UpdateVideo(_ context.Context, id int64, video model.Video) (model.Video, error) {
	for index, current := range r.videos {
		if current.ID == id {
			video.ID = id
			r.videos[index] = video
			return video, nil
		}
	}
	return model.Video{}, dao.ErrVideoNotFound
}

func (r *stubRepository) DeleteVideo(_ context.Context, id int64) error {
	for index, video := range r.videos {
		if video.ID == id {
			r.videos = append(r.videos[:index], r.videos[index+1:]...)
			return nil
		}
	}
	return dao.ErrVideoNotFound
}

func (r *stubRepository) ListTimeline(context.Context) ([]model.TimelineEntry, error) {
	return []model.TimelineEntry{{Period: "2026"}}, nil
}

func (r *stubRepository) CreateComment(_ context.Context, input model.CreateCommentInput) (model.Comment, error) {
	if _, err := r.GetPostBySlug(context.Background(), input.Slug, ""); err != nil {
		return model.Comment{}, err
	}
	comment := model.Comment{
		ID:         int64(len(r.comments[input.Slug]) + 1),
		AuthorName: input.AuthorName,
		Content:    input.Content,
		CreatedAt:  time.Now().UTC(),
	}
	r.comments[input.Slug] = append([]model.Comment{comment}, r.comments[input.Slug]...)
	return comment, nil
}

func (r *stubRepository) ToggleLike(_ context.Context, input model.ToggleLikeInput) (model.LikeState, error) {
	if _, err := r.GetPostBySlug(context.Background(), input.Slug, input.VisitorID); err != nil {
		return model.LikeState{}, err
	}
	if r.likes[input.Slug] == nil {
		r.likes[input.Slug] = make(map[string]bool)
	}
	if r.likes[input.Slug][input.VisitorID] {
		delete(r.likes[input.Slug], input.VisitorID)
		return model.LikeState{LikeCount: len(r.likes[input.Slug]), Liked: false}, nil
	}
	r.likes[input.Slug][input.VisitorID] = true
	return model.LikeState{LikeCount: len(r.likes[input.Slug]), Liked: true}, nil
}
