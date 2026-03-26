package service

import (
	"context"
	"errors"
	"strings"

	"personal/blog/backend/dao"
	"personal/blog/backend/model"
)

var (
	ErrProjectNameRequired    = errors.New("project name is required")
	ErrProjectSummaryRequired = errors.New("project summary is required")
	ErrProjectLinkRequired    = errors.New("project link is required")
)

type AdminProjectService struct {
	repository dao.BlogRepository
}

func NewAdminProjectService(repository dao.BlogRepository) *AdminProjectService {
	return &AdminProjectService{repository: repository}
}

func (s *AdminProjectService) ListProjects(ctx context.Context) ([]model.Project, error) {
	return s.repository.ListProjects(ctx)
}

func (s *AdminProjectService) GetProject(ctx context.Context, id int64) (model.Project, error) {
	project, err := s.repository.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, dao.ErrProjectNotFound) {
			return model.Project{}, dao.ErrProjectNotFound
		}
		return model.Project{}, err
	}
	return project, nil
}

func (s *AdminProjectService) CreateProject(ctx context.Context, input model.CreateProjectInput) (model.Project, error) {
	project, err := normalizeProject(input.Project)
	if err != nil {
		return model.Project{}, err
	}
	return s.repository.CreateProject(ctx, project)
}

func (s *AdminProjectService) UpdateProject(ctx context.Context, input model.UpdateProjectInput) (model.Project, error) {
	project, err := normalizeProject(input.Project)
	if err != nil {
		return model.Project{}, err
	}
	return s.repository.UpdateProject(ctx, input.ID, project)
}

func (s *AdminProjectService) DeleteProject(ctx context.Context, input model.DeleteProjectInput) error {
	return s.repository.DeleteProject(ctx, input.ID)
}

func normalizeProject(project model.Project) (model.Project, error) {
	project.Name = strings.TrimSpace(project.Name)
	project.Summary = strings.TrimSpace(project.Summary)
	project.Status = strings.TrimSpace(project.Status)
	project.Link = strings.TrimSpace(project.Link)
	project.ImageURL = strings.TrimSpace(project.ImageURL)
	project.Accent = strings.TrimSpace(project.Accent)

	if project.Name == "" {
		return model.Project{}, ErrProjectNameRequired
	}
	if project.Summary == "" {
		return model.Project{}, ErrProjectSummaryRequired
	}
	if project.Link == "" {
		return model.Project{}, ErrProjectLinkRequired
	}

	techStack := make([]string, 0, len(project.TechStack))
	for _, item := range project.TechStack {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		techStack = append(techStack, trimmed)
	}
	project.TechStack = techStack
	return project, nil
}
