package domain

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Project defines data related to a project repository
type Project struct {
	gorm.Model    `json:"-"`
	ProjectID     uuid.UUID `gorm:"type:uuid;primary_key;" json:"projectId"`
	Name          string    `json:"name,omitempty"`
	UnzipedPath   string    `json:"unzip,omitempty"`
	ZippedPath    string    `json:"zip,omitempty"`
	GitZippedPath string    `json:"git,omitempty"`
}

// NewProject creates an instance of Project
func NewProject(id uuid.UUID, name, unzipped, zipped, gzp string) *Project {
	return &Project{
		ProjectID:     id,
		Name:          name,
		UnzipedPath:   unzipped,
		ZippedPath:    zipped,
		GitZippedPath: gzp,
	}
}
