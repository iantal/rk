package repository

import (
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/iantal/rk/internal/domain"
	"github.com/jinzhu/gorm"
)

// ProjectDB defines the CRUD operations for storing projects in the db
type ProjectDB struct {
	log hclog.Logger
	db  *gorm.DB
}

// NewProjectDB returns a ProjectDB object for handling CRUD operations
func NewProjectDB(log hclog.Logger, db *gorm.DB) *ProjectDB {
	db.AutoMigrate(&domain.Project{})
	return &ProjectDB{
		log: log,
		db:  db,
	}
}

// AddProject adds a project to the db
func (p *ProjectDB) AddProject(project *domain.Project) {
	p.db.Create(&project)
	return
}

// GetProjects returns all existing projects in the db
func (p *ProjectDB) GetProjects() ([]*domain.Project, error) {
	var projects []*domain.Project
	p.db.Find(&projects)
	return projects, nil
}

// GetProjectByID returns the project with the given id
func (p *ProjectDB) GetProjectByID(id uuid.UUID) *domain.Project {
	return &domain.Project{}
}
