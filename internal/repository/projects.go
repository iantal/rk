package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/iantal/rk/internal/domain"
	"github.com/jinzhu/gorm"
)

type ProjectDB struct {
	log hclog.Logger
	db  *gorm.DB
}

// NewProjectDB returns a ProjectDB object for handling CRUD operations
func NewProjectDB(log hclog.Logger, db *gorm.DB) *ProjectDB {
	return &ProjectDB{
		log: log,
		db: db,
	}
}

func (p *ProjectDB) AddProject() {
	return
}

func (p *ProjectDB) GetProjects() ([]*domain.Project, error) {
	// var projects []*domain.Project
 	// p.db.Find(&projects)
	return []*domain.Project{}, nil
}

func (p *ProjectDB) GetProjectByID(id uuid.UUID) *domain.Project {
	return &domain.Project{}
}
