package handlers

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/iantal/rk/internal/domain"
	"github.com/iantal/rk/internal/files"
	"github.com/iantal/rk/internal/repository"
	"github.com/iantal/rk/internal/util"
	"github.com/sirupsen/logrus"
)

// Projects is a handler for reading and writing projects to a storage and db
type Projects struct {
	l     *util.StandardLogger
	store files.Storage
	db    *repository.ProjectDB
}

// NewProjects creates a handler for projects
func NewProjects(log *util.StandardLogger, store files.Storage, db *repository.ProjectDB) *Projects {
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
		p.l.WithField("error", err).Error("Unable to serialize project")

	}
}

func (p *Projects) Download(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	project := p.db.GetProjectByID(id)
	if project == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		p.l.WithField("projectId", id).Error("Not found")
		util.ToJSON(&GenericError{Message: "Project not found"}, rw)
		return
	}

	rw.Header().Set("Content-type", "application/zip")
	http.ServeFile(rw, r, project.ZippedPath)
}

func (p *Projects) DownloadGitDir(rw http.ResponseWriter, r *http.Request) {
	p.l.Info("Downloading git files only")
	vars := mux.Vars(r)
	id := vars["id"]

	project := p.db.GetProjectByID(id)
	if project == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		p.l.WithField("projectId", id).Error("Not found")
		util.ToJSON(&GenericError{Message: "Project not found"}, rw)
		return
	}

	rw.Header().Set("Content-type", "application/zip")
	http.ServeFile(rw, r, project.GitZippedPath)
}

func (p *Projects) ListSingle(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	project := p.db.GetProjectByID(id)
	if project == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		p.l.WithField("projectId", id).Error("Not found")
		util.ToJSON(&GenericError{Message: "Project not found"}, rw)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	err := util.ToJSON(project, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.WithField("error", err).Error("Unable to serialize project")
	}
}

// CreateProject saves the project zip on the storage and adds it in the database as well
func (p *Projects) CreateProject(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.New()
	fn := vars["filename"]

	p.l.WithFields(logrus.Fields{
		"projectID": id,
		"filename":  fn,
	}).Info("Creating project")

	savedProject := p.save(id, fn, rw, r.Body)
	if savedProject != nil {
		rw.Header().Set("Content-type", "application/json")
		err := util.ToJSON(savedProject, rw)
		if err != nil {
			// we should never be here but log the error just incase
			p.l.WithField("error", err).Error("Unable to serialize project")
		}
	}
}

func (p *Projects) save(id uuid.UUID, path string, rw http.ResponseWriter, r io.ReadCloser) *domain.Project {
	p.l.WithFields(logrus.Fields{
		"projectID": id,
		"path":      path,
	}).Info("Save project to storage")

	unzippedPath := filepath.Join(id.String(), "unzip")
	gzp := filepath.Join(id.String(), "git")

	zp := path + ".zip"
	fp := filepath.Join(id.String(), "zip", zp)
	err := p.store.Save(fp, r)

	if err != nil {
		p.l.WithField("error", err).Error("Unable to save file")
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
		return nil
	}

	go p.process(id.String(), fp, unzippedPath, gzp, path)

	zf := filepath.Join(p.store.FullPath(gzp), zp)
	project := domain.NewProject(id, path, p.store.FullPath(unzippedPath), p.store.FullPath(fp), zf)
	p.db.AddProject(project)
	return project
}

func (p *Projects) process(id, fp, unzippedPath, gzp, path string) {
	p.l.WithFields(logrus.Fields{
		"path": unzippedPath,
		"projectID": id,
	}).Info("Unzipping")

	err := p.store.Unzip(p.store.FullPath(fp), p.store.FullPath(unzippedPath), path)
	if err != nil {
		p.l.WithField("error", err).Error("Unable to unzip file")
	}

	err = p.store.Zip(p.store.FullPath(unzippedPath), p.store.FullPath(gzp), ".git", path)
	if err != nil {
		p.l.WithField("error", err).Error("Unable to zip .git directory")
	}
}
