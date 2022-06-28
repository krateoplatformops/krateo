package storage

import (
	"errors"
	"time"
)

// Errors defined by this package.
var (
	ErrStorageNotFound        = errors.New("storage not found")
	ErrEmptyStorage           = errors.New("storage is empty")
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrEntryNotFound          = errors.New("entry not found")
)

// Entry is a generic representation of a storage item
type Entry struct {
	Meta         Metadata
	Path         string
	Content      []byte
	LastModified time.Time
}

// Metadata represents the meta information of the entry
// includes object name , object version , etc...
type Metadata struct {
	Name    string
	Version string
}

// Storage is a generic interface for storage backends
type Storage interface {
	List(prefix string) ([]Entry, error)
	Get(path string) (Entry, error)
	Put(path string, content []byte) error
	Delete(path string) error
}
