package remote

import (
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
)

type Name interface {
	GetRemoteNameWithUrl() string
}

type remoteName struct {
	repo             middleware.Repository
	remoteUrl        string
	remoteValidation Validation
	remoteList       List
}

func (r remoteName) GetRemoteNameWithUrl() string {
	if validationErr := r.remoteValidation.ValidateRemoteFields(); validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return ""
	}

	remoteList := r.remoteList.GetAllRemotes()
	if remoteList == nil {
		logger.Log("repo has no remotes", global.StatusError)
		return ""
	}

	for _, remote := range remoteList {
		if r.remoteUrl == remote.RemoteURL {
			logger.Log("Matching remote found for the URL", global.StatusInfo)
			return remote.RemoteName
		}
	}

	return ""
}

func NewGetRemoteName(repo middleware.Repository, remoteUrl string, remoteValidation Validation, remoteList List) Name {
	return remoteName{
		repo:             repo,
		remoteUrl:        remoteUrl,
		remoteValidation: remoteValidation,
		remoteList:       remoteList,
	}
}
