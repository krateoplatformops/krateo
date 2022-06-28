package git

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/krateoplatformops/krateo/pkg/storage"
)

const (
	commitAuthorEmail = "krateoctl@krateoplatformops.io"
	commitAuthorName  = "krateoctl"
)

var _ storage.Storage = (*Git)(nil)

// Git is a storage backend for git storage
type Git struct {
	repoUrl  string
	username string
	password string

	storer *memory.Storage
	fs     billy.Filesystem
	repo   *git.Repository
}

type GitOption func(*Git)

func WithGitUser(usr string) GitOption {
	return func(r *Git) {
		r.username = usr
	}
}

func WithGitPassword(pwd string) GitOption {
	return func(r *Git) {
		r.password = pwd
	}
}

func WithGitToken(tkn string) GitOption {
	return func(r *Git) {
		r.username = ""
		r.password = tkn
	}
}

// NewGit creates a new instance of Git repository
func NewGit(repoUrl string, opts ...GitOption) (*Git, error) {
	u, err := url.Parse(repoUrl)
	if err != nil {
		return nil, err
	}

	res := &Git{
		repoUrl: u.String(),
		storer:  memory.NewStorage(),
		fs:      memfs.New(),
	}

	for _, o := range opts {
		o(res)
	}

	res.repo, err = git.Init(res.storer, res.fs)
	if err != nil {
		return nil, err
	}

	_, err = res.repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{res.repoUrl},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Get retrieves an object from root directory
func (s *Git) Get(path string) (storage.Entry, error) {
	err := pull(s)
	if err != nil {
		if errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return storage.Entry{}, storage.ErrEntryNotFound
		}
		return storage.Entry{}, fmt.Errorf("failed to pull from '%s': %w", s.repoUrl, err)
	}

	commit, err := getHeadCommit(s)
	if err != nil {
		return storage.Entry{}, fmt.Errorf("failed to retrieves the branch pointed by HEAD: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return storage.Entry{}, fmt.Errorf("failed to get commit tree: %w", err)
	}

	obj, err := tree.FindEntry(path)
	if err != nil {
		return storage.Entry{}, storage.ErrEntryNotFound
	}

	st, err := s.fs.Stat(path)
	if err != nil {
		return storage.Entry{}, err
	}

	// skip directories
	if st.IsDir() {
		return storage.Entry{}, nil
	}

	res := storage.Entry{
		Path: path,
		Meta: storage.Metadata{
			Name:    obj.Name,
			Version: commit.ID().String(),
		},
		LastModified: st.ModTime(),
	}

	res.Content, err = fetchEntryData(s.fs, path)
	if err != nil {
		return storage.Entry{}, err
	}

	return res, nil
}

func (s *Git) Put(path string, content []byte) error {
	err := pull(s)
	if err != nil {
		if !errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return fmt.Errorf("failed to pull from '%s': %w", s.repoUrl, err)
		}

	}

	fp, err := s.fs.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err := fp.Write(content); err != nil {
		return err
	}

	wt, err := s.repo.Worktree()
	if err != nil {
		return err
	}
	// git add $path
	if _, err := wt.Add(path); err != nil {
		return err
	}

	// git commit -m $message
	_, err = wt.Commit(fmt.Sprintf("add %s", path), &git.CommitOptions{
		Author: &object.Signature{
			Name:  commitAuthorName,
			Email: commitAuthorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	//Push the code to the remote
	return s.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       s.credentials(),
	})
}

// List recursively lists all blobs.
// Returns only entries with the specified prefix eventually.
func (s *Git) List(prefix string) ([]storage.Entry, error) {
	err := pull(s)
	if err != nil {
		if errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return []storage.Entry{}, nil
		}
		return []storage.Entry{}, fmt.Errorf("failed to pull from '%s': %w", s.repoUrl, err)
	}

	commit, err := getHeadCommit(s)
	if err != nil {
		return []storage.Entry{}, fmt.Errorf("failed to retrieves the branch pointed by HEAD: %w", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		return []storage.Entry{}, fmt.Errorf("failed to get commit tree: %w", err)
	}

	return fillEntries(s, tree, commit, prefix)
}

func (s *Git) Delete(path string) error {
	err := pull(s)
	if err != nil {
		if errors.Is(err, transport.ErrEmptyRemoteRepository) {
			return nil
		}
		return fmt.Errorf("failed to pull from '%s': %w", s.repoUrl, err)
	}

	wt, err := s.repo.Worktree()
	if err != nil {
		return err
	}
	// git rm $path
	if _, err := wt.Remove(path); err != nil {
		if errors.Is(err, index.ErrEntryNotFound) {
			return storage.ErrEntryNotFound
		}
		return fmt.Errorf("failed to remove '%s': %w", path, err)
	}

	// git commit -m $message
	_, err = wt.Commit(fmt.Sprintf("remove %s", path), &git.CommitOptions{
		Author: &object.Signature{
			Name:  commitAuthorName,
			Email: commitAuthorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	//Push the code to the remote
	return s.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       s.credentials(),
	})
}

func (s *Git) credentials() *http.BasicAuth {
	if len(s.password) == 0 {
		return nil
	}
	usr := s.username
	if len(usr) == 0 {
		usr = "abc123" // yes, this can be anything except an empty string
	}

	return &http.BasicAuth{
		Username: usr,
		Password: s.password,
	}
}

func pull(s *Git) error {
	// Get the working directory for the repository
	wt, err := s.repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
		//Depth:      1,
		Auth: s.credentials(),
	})

	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			err = nil
		}
	}

	return err
}

func getHeadCommit(s *Git) (*object.Commit, error) {
	// retrieve the branch being pointed by HEAD
	ref, err := s.repo.Head()
	if err != nil {
		return nil, err
	}

	// retrieve the commit object
	return s.repo.CommitObject(ref.Hash())
}

func fillEntries(s *Git, tree *object.Tree, commit *object.Commit, prefix string) ([]storage.Entry, error) {
	entries := make([]storage.Entry, 0)

	seen := make(map[plumbing.Hash]bool)
	iter := object.NewTreeWalker(tree, true, seen)

	withPrefix := len(prefix) > 0

	var name string
	var obj object.TreeEntry
	var err error
	for err == nil {
		name, obj, err = iter.Next()

		if len(name) == 0 {
			continue
		}

		st, err := s.fs.Stat(name)
		if err != nil {
			return nil, err
		}

		if st.IsDir() {
			continue
		}

		data, err := fetchEntryData(s.fs, name)
		if err != nil {
			return nil, err
		}

		entry := storage.Entry{
			Path: name,
			Meta: storage.Metadata{
				Name:    obj.Name,
				Version: commit.ID().String(),
			},
			LastModified: st.ModTime(),
			Content:      data,
		}

		if withPrefix {
			if strings.HasPrefix(name, prefix) {
				entries = append(entries, entry)
			}

			continue
		}

		entries = append(entries, entry)
	}

	if err == io.EOF {
		err = nil
	}

	return entries, err
}

func fetchEntryData(fs billy.Filesystem, path string) ([]byte, error) {
	fp, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return ioutil.ReadAll(fp)
}
