package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"personal/blog/backend/api"
	"personal/blog/backend/dao"
	"personal/blog/backend/middleware"
	"personal/blog/backend/model"
	"personal/blog/backend/service"
)

type AdminController struct {
	authService    *service.AdminAuthService
	contentService *service.AdminContentService
	projectService *service.AdminProjectService
	videoService   *service.AdminVideoService
	cookieName     string
	cookieSecure   bool
}

func NewAdminController(
	authService *service.AdminAuthService,
	contentService *service.AdminContentService,
	projectService *service.AdminProjectService,
	videoService *service.AdminVideoService,
	cookieName string,
	cookieSecure bool,
) *AdminController {
	return &AdminController{
		authService:    authService,
		contentService: contentService,
		projectService: projectService,
		videoService:   videoService,
		cookieName:     cookieName,
		cookieSecure:   cookieSecure,
	}
}

func (c *AdminController) Login(w http.ResponseWriter, r *http.Request) {
	var request api.AdminLoginRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	session, token, err := c.authService.Login(model.AdminLoginInput{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, api.ErrorResponse{Message: err.Error()})
		return
	}

	http.SetCookie(w, c.newSessionCookie(token, session.ExpiresAt))
	writeJSON(w, http.StatusOK, api.AdminLoginResponse{
		Session: api.AdminSessionPayload{
			Authenticated: true,
			Username:      session.Username,
		},
	})
}

func (c *AdminController) Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Secure:   c.cookieSecure,
	})
	writeJSON(w, http.StatusOK, api.AdminLoginResponse{
		Session: api.AdminSessionPayload{
			Authenticated: false,
		},
	})
}

func (c *AdminController) GetSession(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.AdminSessionFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, api.ErrorResponse{Message: "unauthorized"})
		return
	}
	writeJSON(w, http.StatusOK, api.AdminLoginResponse{
		Session: api.AdminSessionPayload{
			Authenticated: true,
			Username:      session.Username,
		},
	})
}

func (c *AdminController) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := c.contentService.ListPosts(r.Context())
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, api.AdminPostsResponse{Posts: toAdminPostSummaries(posts)})
}

func (c *AdminController) HandlePostsRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/posts"), "/")

	switch {
	case path == "" && r.Method == http.MethodGet:
		c.ListPosts(w, r)
	case path == "" && r.Method == http.MethodPost:
		c.CreatePost(w, r)
	case path == "":
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		slug := path
		switch r.Method {
		case http.MethodGet:
			c.GetPost(w, r, slug)
		case http.MethodPut:
			c.UpdatePost(w, r, slug)
		case http.MethodDelete:
			c.DeletePost(w, r, slug)
		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (c *AdminController) HandleVideosRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/videos"), "/")

	switch {
	case path == "thumbnail" && r.Method == http.MethodPost:
		c.SuggestVideoThumbnail(w, r)
	case path == "thumbnail":
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	case path == "" && r.Method == http.MethodGet:
		c.ListVideos(w, r)
	case path == "" && r.Method == http.MethodPost:
		c.CreateVideo(w, r)
	case path == "":
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		id, err := strconv.ParseInt(path, 10, 64)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: "invalid video id"})
			return
		}
		switch r.Method {
		case http.MethodGet:
			c.GetVideo(w, r, id)
		case http.MethodPut:
			c.UpdateVideo(w, r, id)
		case http.MethodDelete:
			c.DeleteVideo(w, r, id)
		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (c *AdminController) HandleProjectsRoute(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/admin/projects"), "/")

	switch {
	case path == "" && r.Method == http.MethodGet:
		c.ListProjectsAdmin(w, r)
	case path == "" && r.Method == http.MethodPost:
		c.CreateProject(w, r)
	case path == "":
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	default:
		id, err := strconv.ParseInt(path, 10, 64)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: "invalid project id"})
			return
		}
		switch r.Method {
		case http.MethodGet:
			c.GetProject(w, r, id)
		case http.MethodPut:
			c.UpdateProjectAdmin(w, r, id)
		case http.MethodDelete:
			c.DeleteProjectAdmin(w, r, id)
		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (c *AdminController) GetPost(w http.ResponseWriter, r *http.Request, slug string) {
	post, err := c.contentService.GetPost(r.Context(), model.GetPostInput{Slug: slug})
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, api.AdminPostResponse{Post: toAdminPostPayload(post)})
}

func (c *AdminController) CreatePost(w http.ResponseWriter, r *http.Request) {
	var request api.AdminSavePostRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	post, err := c.contentService.CreatePost(r.Context(), model.CreatePostInput{
		Post: toModelAdminPost(request),
	})
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, api.AdminPostResponse{Post: toAdminPostPayload(post)})
}

func (c *AdminController) UpdatePost(w http.ResponseWriter, r *http.Request, slug string) {
	var request api.AdminSavePostRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	post, err := c.contentService.UpdatePost(r.Context(), model.UpdatePostInput{
		CurrentSlug: slug,
		Post:        toModelAdminPost(request),
	})
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, api.AdminPostResponse{Post: toAdminPostPayload(post)})
}

func (c *AdminController) DeletePost(w http.ResponseWriter, r *http.Request, slug string) {
	if err := c.contentService.DeletePost(r.Context(), model.DeletePostInput{Slug: slug}); err != nil {
		writeAdminServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *AdminController) ListProjectsAdmin(w http.ResponseWriter, r *http.Request) {
	projects, err := c.projectService.ListProjects(r.Context())
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, api.AdminProjectsResponse{Projects: toAdminProjectPayloads(projects)})
}

func (c *AdminController) GetProject(w http.ResponseWriter, r *http.Request, id int64) {
	project, err := c.projectService.GetProject(r.Context(), id)
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, api.AdminProjectResponse{Project: toAdminProjectPayload(project)})
}

func (c *AdminController) CreateProject(w http.ResponseWriter, r *http.Request) {
	var request api.AdminSaveProjectRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	project, err := c.projectService.CreateProject(r.Context(), model.CreateProjectInput{
		Project: toModelAdminProject(request),
	})
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, api.AdminProjectResponse{Project: toAdminProjectPayload(project)})
}

func (c *AdminController) UpdateProjectAdmin(w http.ResponseWriter, r *http.Request, id int64) {
	var request api.AdminSaveProjectRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	project, err := c.projectService.UpdateProject(r.Context(), model.UpdateProjectInput{
		ID:      id,
		Project: toModelAdminProject(request),
	})
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, api.AdminProjectResponse{Project: toAdminProjectPayload(project)})
}

func (c *AdminController) DeleteProjectAdmin(w http.ResponseWriter, r *http.Request, id int64) {
	if err := c.projectService.DeleteProject(r.Context(), model.DeleteProjectInput{ID: id}); err != nil {
		writeAdminServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *AdminController) ListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := c.videoService.ListVideos(r.Context())
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, api.AdminVideosResponse{Videos: toAdminVideoPayloads(videos)})
}

func (c *AdminController) GetVideo(w http.ResponseWriter, r *http.Request, id int64) {
	video, err := c.videoService.GetVideo(r.Context(), id)
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, api.AdminVideoResponse{Video: toAdminVideoPayload(video)})
}

func (c *AdminController) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var request api.AdminSaveVideoRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	video, err := c.videoService.CreateVideo(r.Context(), model.CreateVideoInput{
		Video: toModelAdminVideo(request),
	})
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, api.AdminVideoResponse{Video: toAdminVideoPayload(video)})
}

func (c *AdminController) UpdateVideo(w http.ResponseWriter, r *http.Request, id int64) {
	var request api.AdminSaveVideoRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	video, err := c.videoService.UpdateVideo(r.Context(), model.UpdateVideoInput{
		ID:    id,
		Video: toModelAdminVideo(request),
	})
	if err != nil {
		writeAdminServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, api.AdminVideoResponse{Video: toAdminVideoPayload(video)})
}

func (c *AdminController) DeleteVideo(w http.ResponseWriter, r *http.Request, id int64) {
	if err := c.videoService.DeleteVideo(r.Context(), model.DeleteVideoInput{ID: id}); err != nil {
		writeAdminServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *AdminController) SuggestVideoThumbnail(w http.ResponseWriter, r *http.Request) {
	var request api.AdminSuggestThumbnailRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	thumbnailURL := c.videoService.SuggestThumbnail(request.URL)
	writeJSON(w, http.StatusOK, api.AdminSuggestThumbnailResponse{ThumbnailURL: thumbnailURL})
}

func toAdminPostSummaries(posts []model.Post) []api.AdminPostSummaryPayload {
	result := make([]api.AdminPostSummaryPayload, 0, len(posts))
	for _, post := range posts {
		result = append(result, api.AdminPostSummaryPayload{
			Slug:         post.Slug,
			Title:        post.Title,
			Summary:      post.Summary,
			Category:     post.Category,
			ReadTime:     post.ReadTime,
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

func toAdminPostPayload(post model.Post) api.AdminPostPayload {
	blocks := make([]api.AdminContentBlock, 0, len(post.Blocks))
	for _, block := range post.Blocks {
		blocks = append(blocks, api.AdminContentBlock{
			Kind:  block.Kind,
			Title: block.Title,
			Text:  block.Text,
			URL:   block.URL,
			Items: block.Items,
		})
	}

	return api.AdminPostPayload{
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
		Blocks:       blocks,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
	}
}

func toModelAdminPost(request api.AdminSavePostRequest) model.Post {
	blocks := make([]model.ContentBlock, 0, len(request.Blocks))
	for _, block := range request.Blocks {
		blocks = append(blocks, model.ContentBlock{
			Kind:  block.Kind,
			Title: block.Title,
			Text:  block.Text,
			URL:   block.URL,
			Items: block.Items,
		})
	}

	return model.Post{
		Slug:        request.Slug,
		Title:       request.Title,
		Summary:     request.Summary,
		Category:    request.Category,
		ReadTime:    request.ReadTime,
		HeroNote:    request.HeroNote,
		CoverLabel:  request.CoverLabel,
		Tags:        request.Tags,
		Featured:    request.Featured,
		PublishedAt: parseAdminDate(request.PublishedAt),
		Blocks:      blocks,
	}
}

func toAdminVideoPayloads(videos []model.Video) []api.AdminVideoPayload {
	result := make([]api.AdminVideoPayload, 0, len(videos))
	for _, video := range videos {
		result = append(result, toAdminVideoPayload(video))
	}
	return result
}

func toAdminProjectPayloads(projects []model.Project) []api.AdminProjectPayload {
	result := make([]api.AdminProjectPayload, 0, len(projects))
	for _, project := range projects {
		result = append(result, toAdminProjectPayload(project))
	}
	return result
}

func toAdminProjectPayload(project model.Project) api.AdminProjectPayload {
	return api.AdminProjectPayload{
		ID:        project.ID,
		Name:      project.Name,
		Summary:   project.Summary,
		Status:    project.Status,
		Link:      project.Link,
		ImageURL:  project.ImageURL,
		Accent:    project.Accent,
		TechStack: project.TechStack,
	}
}

func toModelAdminProject(request api.AdminSaveProjectRequest) model.Project {
	return model.Project{
		Name:      request.Name,
		Summary:   request.Summary,
		Status:    request.Status,
		Link:      request.Link,
		ImageURL:  request.ImageURL,
		Accent:    request.Accent,
		TechStack: request.TechStack,
	}
}

func toAdminVideoPayload(video model.Video) api.AdminVideoPayload {
	return api.AdminVideoPayload{
		ID:          video.ID,
		Title:       video.Title,
		Description: video.Description,
		URL:         video.URL,
		ThumbnailURL: video.ThumbnailURL,
		PublishedAt: video.PublishedAt.Format("2006-01-02"),
	}
}

func toModelAdminVideo(request api.AdminSaveVideoRequest) model.Video {
	return model.Video{
		Title:       request.Title,
		Description: request.Description,
		URL:         request.URL,
		ThumbnailURL: request.ThumbnailURL,
		PublishedAt: parseAdminDate(request.PublishedAt),
	}
}

func parseAdminDate(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}

	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return parsed
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed
	}
	return time.Time{}
}

func writeAdminServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrPostNotFound):
		writeJSON(w, http.StatusNotFound, api.ErrorResponse{Message: err.Error()})
	case errors.Is(err, service.ErrPostSlugRequired),
		errors.Is(err, service.ErrPostSlugInvalid),
		errors.Is(err, service.ErrPostTitleRequired),
		errors.Is(err, service.ErrPostSummaryRequired),
		errors.Is(err, service.ErrPostCategoryRequired),
		errors.Is(err, service.ErrPostReadTimeRequired),
		errors.Is(err, service.ErrPostPublishedAtInvalid),
		errors.Is(err, service.ErrPostBlocksRequired),
		errors.Is(err, service.ErrDuplicatePostSlug),
		errors.Is(err, service.ErrProjectNameRequired),
		errors.Is(err, service.ErrProjectSummaryRequired),
		errors.Is(err, service.ErrProjectLinkRequired),
		errors.Is(err, service.ErrVideoTitleRequired),
		errors.Is(err, service.ErrVideoURLRequired),
		errors.Is(err, service.ErrVideoPublishedAtInvalid):
		writeJSON(w, http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
	case errors.Is(err, dao.ErrProjectNotFound):
		writeJSON(w, http.StatusNotFound, api.ErrorResponse{Message: err.Error()})
	case errors.Is(err, dao.ErrVideoNotFound):
		writeJSON(w, http.StatusNotFound, api.ErrorResponse{Message: err.Error()})
	default:
		writeInternalError(w, err)
	}
}

func (c *AdminController) newSessionCookie(token string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     c.cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiresAt,
		SameSite: http.SameSiteLaxMode,
		Secure:   c.cookieSecure,
	}
}
