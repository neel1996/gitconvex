package remote

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type ListRemoteUrl interface {
	GetAllRemoteUrl() []*string
}

type listRemoteUrl struct {
	repo       *git2go.Repository
	validation Validation
}

func (u listRemoteUrl) GetAllRemoteUrl() []*string {
	var remoteURL []*string

	repo := u.repo

	if validationErr := u.validation.ValidateRemoteFields(repo); validationErr != nil {
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

func NewRemoteUrlData(repo *git2go.Repository, validation Validation) ListRemoteUrl {
	return listRemoteUrl{repo: repo, validation: validation}
}
