package local

import (
	"io/ioutil"
	"os"
	pathutil "path"
	"path/filepath"

	"github.com/krateoplatformops/krateo/pkg/storage"
)

var _ storage.Storage = (*LocalFilesystem)(nil)

// LocalFilesystem is a storage backend for local filesystem storage
type LocalFilesystem struct {
	RootDir string
}

// NewLocalFilesystem creates a new instance of LocalFilesystemBackend
func NewLocalFilesystem(root string) *LocalFilesystem {
	absPath, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	b := &LocalFilesystem{RootDir: absPath}
	return b
}

// List lists all objects in root directory (depth 1)
func (b *LocalFilesystem) List(prefix string) ([]storage.Entry, error) {
	var objects []storage.Entry
	files, err := ioutil.ReadDir(pathutil.Join(b.RootDir, prefix))
	if err != nil {
		if os.IsNotExist(err) { // OK if the directory doesnt exist yet
			err = nil
		}
		return objects, err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		object := storage.Entry{Path: f.Name(), Content: []byte{}, LastModified: f.ModTime()}
		objects = append(objects, object)
	}
	return objects, nil
}

// Get retrieves an object from root directory
func (b *LocalFilesystem) Get(path string) (storage.Entry, error) {
	var object storage.Entry
	object.Path = path
	fullpath := pathutil.Join(b.RootDir, path)
	content, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return object, err
	}
	object.Content = content
	info, err := os.Stat(fullpath)
	if err != nil {
		return object, err
	}
	object.LastModified = info.ModTime()
	return object, err
}

// PutObject puts an object in root directory
func (b *LocalFilesystem) Put(path string, content []byte) error {
	fullpath := pathutil.Join(b.RootDir, path)
	folderPath := pathutil.Dir(fullpath)
	_, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(folderPath, 0774)
			if err != nil {
				return err
			}
			// os.MkdirAll set the dir permissions before the umask
			// we need to use os.Chmod to ensure the permissions of the created directory are 774
			// because the default umask will prevent that and cause the permissions to be 755
			err = os.Chmod(folderPath, 0774)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	err = ioutil.WriteFile(fullpath, content, 0644)
	return err
}

// DeleteObject removes an object from root directory
func (b *LocalFilesystem) Delete(path string) error {
	fullpath := pathutil.Join(b.RootDir, path)
	err := os.Remove(fullpath)
	return err
}
