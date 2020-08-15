package handlers

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/iantal/rk/internal/files"
	"github.com/iantal/rk/internal/repository"
	"github.com/iantal/rk/internal/util"
)

// Projects is a handler for reading and writing projects to a storage and db
type Projects struct {
	l     hclog.Logger
	store files.Storage
	db    *repository.ProjectDB
}

// NewProjects creates a handler for projects
func NewProjects(log hclog.Logger, store files.Storage, db *repository.ProjectDB) *Projects {
	return &Projects{
		l:     log,
		store: store,
		db:    db,
	}
}

// GenericError represents an error of the system
type GenericError struct {
	Message string `json:"message"`
}

// ListAll returns a list of existing projects
func (p *Projects) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Debug("Get all projects")
	rw.Header().Add("Content-Type", "application/json")

	projects, err := p.db.GetProjects()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		util.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = util.ToJSON(projects, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable to serialize project", "error", err)
	}
}
