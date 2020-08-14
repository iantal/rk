package domain

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Project defines data related to a project repository
type Project struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string
	Path string
}
