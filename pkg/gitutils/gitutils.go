package gitutils

import (
	"io/ioutil"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

func ReadFile(fs billy.Filesystem, fn string) ([]byte, error) {
	fp, err := fs.Open(fn)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return ioutil.ReadAll(fp)
}

func Clone(url, token string) (billy.Filesystem, error) {
	fs := memfs.New()

	opts := git.CloneOptions{
		URL:   url,
		Depth: 1,
	}

	if len(token) > 0 {
		opts.Auth = &http.BasicAuth{
			Username: "krateo",
			Password: token,
		}
	}

	_, err := git.Clone(memory.NewStorage(), fs, &opts)
	if err != nil {
		return nil, err
	}

	return fs, nil
}
