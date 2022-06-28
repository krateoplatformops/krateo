package tools

import (
	"fmt"
	"net/url"
	"path"

	"github.com/krateoplatformops/krateo/pkg/storage"
	"github.com/krateoplatformops/krateo/pkg/storage/git"
)

const (
	krateopGitUrl = "https://github.com/krateoplatformops"
)

func GitPullModuleDefinition(name string) (storage.Entry, error) {
	return GitPullModuleDefinitionFromUserRepo(name, krateopGitUrl)
}

func GitPullModuleDefinitionFromUserRepo(name string, gitUrl string, opts ...git.GitOption) (storage.Entry, error) {
	repoUrl, err := url.Parse(gitUrl)
	if err != nil {
		return storage.Entry{}, err
	}
	if gitUrl == krateopGitUrl {
		repoUrl.Path = path.Join(repoUrl.Path, fmt.Sprintf("krateo-module-%s", name))
	}

	store, err := git.NewGit(repoUrl.String(), opts...)
	if err != nil {
		return storage.Entry{}, err
	}

	return store.Get("cluster/definition.yaml")
}

func GitPullModulePackage(name string) (storage.Entry, error) {
	return GitPullModulePackageFromUserRepo(name, krateopGitUrl)
}

func GitPullModulePackageFromUserRepo(name string, gitUrl string, opts ...git.GitOption) (storage.Entry, error) {
	repoUrl, err := url.Parse(gitUrl)
	if err != nil {
		return storage.Entry{}, err
	}
	if gitUrl == krateopGitUrl {
		repoUrl.Path = path.Join(repoUrl.Path, fmt.Sprintf("krateo-module-%s", name))
	}

	store, err := git.NewGit(repoUrl.String(), opts...)
	if err != nil {
		return storage.Entry{}, err
	}

	return store.Get(fmt.Sprintf("defaults/krateo-package-module-%s.yaml", name))
}

func GitPullModuleDefaults(name string) (storage.Entry, error) {
	return GitPullModuleDefaultsFormUserRepo(name, krateopGitUrl)
}

func GitPullModuleDefaultsFormUserRepo(name string, gitUrl string, opts ...git.GitOption) (storage.Entry, error) {
	repoUrl, err := url.Parse(gitUrl)
	if err != nil {
		return storage.Entry{}, err
	}
	if gitUrl == krateopGitUrl {
		repoUrl.Path = path.Join(repoUrl.Path, fmt.Sprintf("krateo-module-%s", name))
	}

	store, err := git.NewGit(repoUrl.String(), opts...)
	if err != nil {
		return storage.Entry{}, err
	}

	return store.Get(fmt.Sprintf("defaults/krateo-module-%s.yaml", name))
}

func GitPushEntry(gitRemoteUrl string, entry storage.Entry, opts ...git.GitOption) error {
	repoUrl, err := url.Parse(gitRemoteUrl)
	if err != nil {
		return err
	}

	git, err := git.NewGit(repoUrl.String(), opts...)
	if err != nil {
		return err
	}

	return git.Put(entry.Path, entry.Content)
}
