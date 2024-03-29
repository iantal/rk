package repository

import (
	"github.com/google/uuid"
	"github.com/iantal/rk/internal/domain"
	"github.com/iantal/rk/internal/util"
	"github.com/jinzhu/gorm"
)

// ProjectDB defines the CRUD operations for storing projects in the db
type ProjectDB struct {
	log *util.StandardLogger
	db  *gorm.DB
}

// NewProjectDB returns a ProjectDB object for handling CRUD operations
func NewProjectDB(log *util.StandardLogger, db *gorm.DB) *ProjectDB {
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
func (p *ProjectDB) GetProjectByID(id string) *domain.Project {
	project := &domain.Project{}
	uid, err := uuid.Parse(id)
	if err != nil {
		p.log.Error("Project with projectId %s was not found", id)
		return nil
	}
	p.db.Find(&project, "project_id = ?", uid)
	return project
}
