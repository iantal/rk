package handlers

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/iantal/rk/internal/domain"
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

// CreateProject saves the project zip on the storage and adds it in the database as well
func (p *Projects) CreateProject(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.New()
	fn := vars["filename"]

	p.l.Info("Creating project", "id", id, "filename", fn)

	// no need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters
	p.save(id, fn, rw, r.Body)

}

func (p *Projects) save(id uuid.UUID, path string, rw http.ResponseWriter, r io.ReadCloser) {
	p.l.Info("Save project - storage", "id", id, "path", path)

	fp := filepath.Join(id.String(), path)
	err := p.store.Save(fp, r)
	if err != nil {
		p.l.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	} else {
		project := domain.NewProject(id, path, p.store.FullPath(fp))
		p.l.Debug("Save project - db", "id", id, "path", path)
		p.db.AddProject(project)
	}
}
