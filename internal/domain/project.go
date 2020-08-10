package domain

import "github.com/google/uuid"

// Project defines data related to a project repository
type Project struct {
	ID   uuid.UUID
	Name string
	Path string
}
