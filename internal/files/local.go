package files

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// Local is an implementation of the Storage interface which works with the
// local disk on the current machine
type Local struct {
	maxFileSize int // maximum numbber of bytes for files
	basePath    string
}

// NewLocal creates a new Local filesytem with the given base path
// basePath is the base directory to save files to
// maxSize is the max number of bytes that a file can be
func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: p}, nil
}

// Save the contents of the Writer to the given path
// path is a relative path, basePath will be appended
func (l *Local) Save(path string, contents io.Reader) error {
	// get the full path for the file
	fp := l.FullPath(path)

	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	// if the file exists delete it
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		return xerrors.Errorf("Unable to get file info: %w", err)
	}

	// create a new file at the path
	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	// write the contents to the new file
	// ensure that we are not writing greater than max bytes
	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to file: %w", err)
	}

	return nil
}

// Get the file at the given path and return a Reader
// the calling function is responsible for closing the reader
func (l *Local) Get(path string) (*os.File, error) {
	// get the full path for the file
	fp := l.FullPath(path)

	// open the file
	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("Unable to open file: %w", err)
	}

	return f, nil
}

// FullPath returns the absolute path
func (l *Local) FullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}

func (l *Local) Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return xerrors.Errorf("Unable to open zip file: %w", err)
	}

	td := filepath.Join(target, "kafka")
	if err := os.MkdirAll(td, 0755); err != nil {
		return xerrors.Errorf("Unable to create target directory: %w", err)
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			fileReader, err := file.Open()
			if err != nil {
				return xerrors.Errorf("Unable to open file: %w", err)
			}
			defer fileReader.Close()

			// create a new file at the path
			targetFile, err := os.Create(path)
			if err != nil {
				return xerrors.Errorf("Unable to create file: %w", err)
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, fileReader); err != nil {
				return xerrors.Errorf("Unable to write target file: %w", err)
			}
		}
	}

	return nil
}
