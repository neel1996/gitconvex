package remote

import (
	"errors"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type ListRemoteUrl interface {
	GetAllRemoteUrl() []*string
}

type listRemoteUrl struct {
	repo *git2go.Repository
}

func (u listRemoteUrl) GetAllRemoteUrl() []*string {
	var remoteURL []*string

	repo := u.repo

	if validationErr := u.validateFields(repo); validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return nil
	}

	remotes, listErr := repo.Remotes.List()
	if listErr != nil {
		logger.Log(listErr.Error(), global.StatusError)
		return nil
	}

	for _, remoteName := range remotes {
		remote, _ := repo.Remotes.Lookup(remoteName)
		if remote != nil {
			url := remote.Url()
			remoteURL = append(remoteURL, &url)
		}
	}

	if len(remoteURL) == 0 {
		logger.Log("No remotes present in the repo", global.StatusError)
		return nil
	}

	return remoteURL
}

func (u listRemoteUrl) validateFields(repo *git2go.Repository) error {
	if repo == nil {
		return errors.New("repo is nil")
	}

	remoteCollection := repo.Remotes
	if remoteCollection == (git2go.RemoteCollection{}) {
		return errors.New("remote collection is empty")
	}

	return nil
}

func NewRemoteUrlData(repo *git2go.Repository) ListRemoteUrl {
	return listRemoteUrl{repo: repo}
}
