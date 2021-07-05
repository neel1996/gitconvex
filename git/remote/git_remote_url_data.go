package remote

import (
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

	if validationErr := NewRemoteValidation(repo).ValidateRemoteFields(); validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return nil
	}

	remoteList := NewRemoteList(u.repo).GetAllRemotes()
	if remoteList == nil {
		logger.Log("repo has no remotes", global.StatusError)
		return nil
	}

	for _, remote := range remoteList {
		remoteURL = append(remoteURL, &remote.RemoteURL)
	}

	if len(remoteURL) == 0 {
		logger.Log("No remotes present in the repo", global.StatusError)
		return nil
	}

	return remoteURL
}

func NewRemoteUrlData(repo *git2go.Repository) ListRemoteUrl {
	return listRemoteUrl{repo: repo}
}
