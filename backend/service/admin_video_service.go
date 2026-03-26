package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"personal/blog/backend/dao"
	"personal/blog/backend/model"
)

var (
	ErrVideoTitleRequired      = errors.New("video title is required")
	ErrVideoURLRequired        = errors.New("video url is required")
	ErrVideoPublishedAtInvalid = errors.New("video published date is invalid")
)

var bilibiliBVIDPattern = regexp.MustCompile(`BV[0-9A-Za-z]+`)

type AdminVideoService struct {
	repository dao.BlogRepository
}

func NewAdminVideoService(repository dao.BlogRepository) *AdminVideoService {
	return &AdminVideoService{repository: repository}
}

func (s *AdminVideoService) ListVideos(ctx context.Context) ([]model.Video, error) {
	return s.repository.ListVideos(ctx)
}

func (s *AdminVideoService) GetVideo(ctx context.Context, id int64) (model.Video, error) {
	video, err := s.repository.GetVideoByID(ctx, id)
	if err != nil {
		if errors.Is(err, dao.ErrVideoNotFound) {
			return model.Video{}, dao.ErrVideoNotFound
		}
		return model.Video{}, err
	}
	return video, nil
}

func (s *AdminVideoService) CreateVideo(ctx context.Context, input model.CreateVideoInput) (model.Video, error) {
	video, err := normalizeVideo(input.Video)
	if err != nil {
		return model.Video{}, err
	}
	return s.repository.CreateVideo(ctx, video)
}

func (s *AdminVideoService) UpdateVideo(ctx context.Context, input model.UpdateVideoInput) (model.Video, error) {
	video, err := normalizeVideo(input.Video)
	if err != nil {
		return model.Video{}, err
	}
	updated, err := s.repository.UpdateVideo(ctx, input.ID, video)
	if err != nil {
		if errors.Is(err, dao.ErrVideoNotFound) {
			return model.Video{}, dao.ErrVideoNotFound
		}
		return model.Video{}, err
	}
	return updated, nil
}

func (s *AdminVideoService) DeleteVideo(ctx context.Context, input model.DeleteVideoInput) error {
	err := s.repository.DeleteVideo(ctx, input.ID)
	if err != nil {
		if errors.Is(err, dao.ErrVideoNotFound) {
			return dao.ErrVideoNotFound
		}
		return err
	}
	return nil
}

func (s *AdminVideoService) SuggestThumbnail(videoURL string) string {
	return deriveVideoThumbnail(strings.TrimSpace(videoURL))
}

func normalizeVideo(video model.Video) (model.Video, error) {
	video.Title = strings.TrimSpace(video.Title)
	video.Description = strings.TrimSpace(video.Description)
	video.URL = strings.TrimSpace(video.URL)
	video.ThumbnailURL = strings.TrimSpace(video.ThumbnailURL)

	if video.Title == "" {
		return model.Video{}, ErrVideoTitleRequired
	}
	if video.URL == "" {
		return model.Video{}, ErrVideoURLRequired
	}
	if video.PublishedAt.IsZero() {
		return model.Video{}, ErrVideoPublishedAtInvalid
	}

	if video.ThumbnailURL == "" {
		video.ThumbnailURL = deriveVideoThumbnail(video.URL)
	}

	return video, nil
}

func deriveVideoThumbnail(videoURL string) string {
	if thumbnail := deriveYouTubeThumbnail(videoURL); thumbnail != "" {
		return thumbnail
	}
	if thumbnail := deriveBilibiliThumbnail(videoURL); thumbnail != "" {
		return thumbnail
	}
	return ""
}

func deriveYouTubeThumbnail(videoURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(videoURL))
	if err != nil {
		return ""
	}

	switch {
	case strings.Contains(parsed.Hostname(), "youtube.com"):
		videoID := parsed.Query().Get("v")
		if videoID != "" {
			return "https://i.ytimg.com/vi/" + videoID + "/hqdefault.jpg"
		}
	case strings.Contains(parsed.Hostname(), "youtu.be"):
		videoID := strings.Trim(parsed.Path, "/")
		if videoID != "" {
			return "https://i.ytimg.com/vi/" + videoID + "/hqdefault.jpg"
		}
	}
	return ""
}

func deriveBilibiliThumbnail(videoURL string) string {
	bvid := bilibiliBVIDPattern.FindString(videoURL)
	if bvid == "" {
		return ""
	}

	requestURL := "https://api.bilibili.com/x/web-interface/view?bvid=" + bvid
	client := &http.Client{Timeout: 4 * time.Second}
	response, err := client.Get(requestURL)
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	var payload struct {
		Code int `json:"code"`
		Data struct {
			Pic string `json:"pic"`
		} `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return ""
	}
	if payload.Code != 0 {
		return ""
	}
	if payload.Data.Pic == "" || strings.Contains(payload.Data.Pic, "transparent.png") {
		return ""
	}
	return payload.Data.Pic
}
