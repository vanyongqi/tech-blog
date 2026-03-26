package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"personal/blog/backend/api"
	"personal/blog/backend/middleware"
	"personal/blog/backend/model"
	"personal/blog/backend/service"
)

type BlogController struct {
	blogService *service.BlogService
}

func NewBlogController(blogService *service.BlogService) *BlogController {
	return &BlogController{blogService: blogService}
}

func (c *BlogController) GetHome(w http.ResponseWriter, r *http.Request) {
	site, err := c.blogService.GetSiteProfile(r.Context())
	if err != nil {
		writeInternalError(w, err)
		return
	}
	featuredPosts, err := c.blogService.ListPosts(r.Context(), model.ListPostsInput{FeaturedOnly: true, Limit: 2})
	if err != nil {
		writeInternalError(w, err)
		return
	}
	recentPosts, err := c.blogService.ListPosts(r.Context(), model.ListPostsInput{Limit: 3})
	if err != nil {
		writeInternalError(w, err)
		return
	}
	projects, err := c.blogService.ListProjects(r.Context())
	if err != nil {
		writeInternalError(w, err)
		return
	}
	videos, err := c.blogService.ListVideos(r.Context())
	if err != nil {
		writeInternalError(w, err)
		return
	}
	timeline, err := c.blogService.ListTimeline(r.Context())
	if err != nil {
		writeInternalError(w, err)
		return
	}

	response := api.HomeResponse{
		Site:          toSitePayload(site),
		FeaturedPosts: toPostSummaries(featuredPosts),
		RecentPosts:   toPostSummaries(recentPosts),
		Projects:      toProjectPayloads(projects),
		Videos:        toVideoPayloads(videos),
		Timeline:      toTimelinePayloads(timeline),
	}
	writeJSON(w, http.StatusOK, response)
}

func (c *BlogController) ListPosts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))

	posts, err := c.blogService.ListPosts(r.Context(), model.ListPostsInput{
		FeaturedOnly: strings.EqualFold(query.Get("featured"), "true"),
		Tag:          query.Get("tag"),
		Limit:        limit,
	})
	if err != nil {
		writeInternalError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, api.PostsResponse{Posts: toPostSummaries(posts)})
}

func (c *BlogController) HandlePostRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/posts/"), "/")
	if path == "" {
		writeJSON(w, http.StatusNotFound, api.ErrorResponse{Message: "post route not found"})
		return
	}

	parts := strings.Split(path, "/")
	slug := parts[0]
	switch {
	case len(parts) == 1 && r.Method == http.MethodGet:
		c.getPost(w, r, slug)
	case len(parts) == 2 && parts[1] == "comments" && r.Method == http.MethodPost:
		c.createComment(w, r, slug)
	case len(parts) == 2 && parts[1] == "likes" && r.Method == http.MethodPost:
		c.toggleLike(w, r, slug)
	case len(parts) == 1:
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
	case len(parts) == 2 && (parts[1] == "comments" || parts[1] == "likes"):
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		writeJSON(w, http.StatusNotFound, api.ErrorResponse{Message: "post route not found"})
	}
}

func (c *BlogController) getPost(w http.ResponseWriter, r *http.Request, slug string) {
	visitor := middleware.MustVisitorIdentity(r.Context())
	post, err := c.blogService.GetPost(r.Context(), model.GetPostInput{
		Slug:      slug,
		VisitorID: visitor.ID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, api.PostResponse{
		Post:    toPostDetailPayload(post),
		Visitor: api.VisitorPayload{DisplayName: visitor.DisplayName},
	})
}

func (c *BlogController) createComment(w http.ResponseWriter, r *http.Request, slug string) {
	visitor := middleware.MustVisitorIdentity(r.Context())
	var request api.CreateCommentRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	comment, err := c.blogService.CreateComment(r.Context(), model.CreateCommentInput{
		Slug:       slug,
		AuthorName: visitor.DisplayName,
		Content:    request.Content,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	post, err := c.blogService.GetPost(r.Context(), model.GetPostInput{
		Slug:      slug,
		VisitorID: visitor.ID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, api.CommentResponse{
		Comment:      toCommentPayload(comment),
		CommentCount: post.CommentCount,
		Visitor:      api.VisitorPayload{DisplayName: visitor.DisplayName},
	})
}

func (c *BlogController) toggleLike(w http.ResponseWriter, r *http.Request, slug string) {
	visitor := middleware.MustVisitorIdentity(r.Context())
	likeState, err := c.blogService.ToggleLike(r.Context(), model.ToggleLikeInput{
		Slug:      slug,
		VisitorID: visitor.ID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, api.LikeResponse{
		LikeCount: likeState.LikeCount,
		Liked:     likeState.Liked,
		Visitor:   api.VisitorPayload{DisplayName: visitor.DisplayName},
	})
}

func toSitePayload(site model.SiteProfile) api.SitePayload {
	stats := make([]api.SiteStat, 0, len(site.Stats))
	for _, stat := range site.Stats {
		stats = append(stats, api.SiteStat{Label: stat.Label, Value: stat.Value})
	}

	links := make([]api.SocialLink, 0, len(site.SocialLinks))
	for _, link := range site.SocialLinks {
		links = append(links, api.SocialLink{Label: link.Label, URL: link.URL})
	}

	return api.SitePayload{
		Name:        site.Name,
		Headline:    site.Headline,
		Intro:       site.Intro,
		Location:    site.Location,
		Domain:      site.Domain,
		Email:       site.Email,
		Motto:       site.Motto,
		TechStack:   site.TechStack,
		Stats:       stats,
		SocialLinks: links,
	}
}

func toPostSummaries(posts []model.Post) []api.PostSummaryPayload {
	result := make([]api.PostSummaryPayload, 0, len(posts))
	for _, post := range posts {
		result = append(result, api.PostSummaryPayload{
			Slug:         post.Slug,
			Title:        post.Title,
			Summary:      post.Summary,
			Category:     post.Category,
			ReadTime:     post.ReadTime,
			HeroNote:     post.HeroNote,
			CoverLabel:   post.CoverLabel,
			Tags:         post.Tags,
			Featured:     post.Featured,
			PublishedAt:  post.PublishedAt.Format("2006-01-02"),
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
		})
	}
	return result
}

func toPostDetailPayload(post model.Post) api.PostDetailPayload {
	blocks := make([]api.ContentBlock, 0, len(post.Blocks))
	for _, block := range post.Blocks {
		blocks = append(blocks, api.ContentBlock{
			Kind:  block.Kind,
			Title: block.Title,
			Text:  block.Text,
			URL:   block.URL,
			Items: block.Items,
		})
	}

	comments := make([]api.CommentPayload, 0, len(post.Comments))
	for _, comment := range post.Comments {
		comments = append(comments, toCommentPayload(comment))
	}

	return api.PostDetailPayload{
		PostSummaryPayload: toPostSummaries([]model.Post{post})[0],
		Blocks:             blocks,
		LikedByVisitor:     post.LikedByVisitor,
		Comments:           comments,
	}
}

func toProjectPayloads(projects []model.Project) []api.ProjectPayload {
	result := make([]api.ProjectPayload, 0, len(projects))
	for _, project := range projects {
		result = append(result, api.ProjectPayload{
			Name:      project.Name,
			Summary:   project.Summary,
			Status:    project.Status,
			Link:      project.Link,
			ImageURL:  project.ImageURL,
			Accent:    project.Accent,
			TechStack: project.TechStack,
		})
	}
	return result
}

func toVideoPayloads(videos []model.Video) []api.VideoPayload {
	result := make([]api.VideoPayload, 0, len(videos))
	for _, video := range videos {
		result = append(result, api.VideoPayload{
			ID:          video.ID,
			Title:       video.Title,
			Description: video.Description,
			URL:         video.URL,
			ThumbnailURL: video.ThumbnailURL,
			PublishedAt: video.PublishedAt.Format("2006-01-02"),
		})
	}
	return result
}

func toTimelinePayloads(entries []model.TimelineEntry) []api.TimelineEntryPayload {
	result := make([]api.TimelineEntryPayload, 0, len(entries))
	for _, entry := range entries {
		result = append(result, api.TimelineEntryPayload{
			Period:      entry.Period,
			Title:       entry.Title,
			Description: entry.Description,
		})
	}
	return result
}

func toCommentPayload(comment model.Comment) api.CommentPayload {
	return api.CommentPayload{
		ID:         comment.ID,
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		CreatedAt:  comment.CreatedAt.Format(time.RFC3339),
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.New("invalid request body")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("invalid request body")
	}
	return nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrPostNotFound):
		writeJSON(w, http.StatusNotFound, api.ErrorResponse{Message: err.Error()})
	case errors.Is(err, service.ErrCommentContentRequired),
		errors.Is(err, service.ErrCommentContentTooLong),
		errors.Is(err, service.ErrVisitorIdentityRequired):
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
	default:
		writeInternalError(w, err)
	}
}

func writeInternalError(w http.ResponseWriter, _ error) {
	writeJSON(w, http.StatusInternalServerError, api.ErrorResponse{Message: "internal server error"})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
