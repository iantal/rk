package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/iantal/rk/internal/domain"
)

type ProjectDB struct {
	log hclog.Logger
}

func (p *ProjectDB) AddProject() {
	return
}

func (p *ProjectDB) GetProjects() []*domain.Project {
	return []*domain.Project{}
}

func (p *ProjectDB) GetProjectByID(id uuid.UUID) *domain.Project {
	return &domain.Project{}
}
