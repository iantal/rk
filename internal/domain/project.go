package domain

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Project defines data related to a project repository
type Project struct {
	gorm.Model `json:"-"`
	ProjectID uuid.UUID `gorm:"type:uuid;primary_key;" json:"projectId"`
	Name      string `json:"name"`
	Path      string `json:"path"`
}

// NewProject creates an instance of Project
func NewProject(id uuid.UUID, name, path string) *Project {
	return &Project{
		ProjectID: id,
		Name:      name,
		Path:      path,
	}
}
