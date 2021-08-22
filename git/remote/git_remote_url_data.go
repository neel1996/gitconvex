package remote

import (
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
)

type ListRemoteUrl interface {
	GetAllRemoteUrl() []*string
}

type listRemoteUrl struct {
	repo             middleware.Repository
	remoteValidation Validation
	remoteList       List
}

func (u listRemoteUrl) GetAllRemoteUrl() []*string {
	var remoteURL []*string

	if validationErr := u.remoteValidation.ValidateRemoteFields(); validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return nil
	}

	remoteList := u.remoteList.GetAllRemotes()
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

func NewRemoteUrlData(repo middleware.Repository, remoteValidation Validation, remoteList List) ListRemoteUrl {
	return listRemoteUrl{
		repo:             repo,
		remoteValidation: remoteValidation,
		remoteList:       remoteList,
	}
}
