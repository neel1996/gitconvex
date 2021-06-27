package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type Name interface {
	GetRemoteNameWithUrl() string
}

type remoteName struct {
	repo      *git2go.Repository
	remoteUrl string
	validate  Validation
}

func (r remoteName) GetRemoteNameWithUrl() string {
	repo := r.repo

	if validationErr := r.validate.ValidateRemoteFields(repo); validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return ""
	}

	remoteList, remoteListError := repo.Remotes.List()
	if remoteListError != nil {
		logger.Log(remoteListError.Error(), global.StatusError)
		return ""
	}

	for _, remote := range remoteList {
		remoteEntry, remoteLookupErr := repo.Remotes.Lookup(remote)
		if remoteLookupErr != nil {
			logger.Log(remoteLookupErr.Error(), global.StatusError)
			continue
		}

		isRemotePresent := r.isRemoteAvailable(remoteEntry, remote)
		if isRemotePresent {
			logger.Log("Matching remote found for the URL", global.StatusInfo)
			return remote
		}
	}

	return ""
}

func (r remoteName) isRemoteAvailable(remoteEntry *git2go.Remote, remote string) bool {
	if remoteEntry.Url() == r.remoteUrl {
		logger.Log(fmt.Sprintf("Remote - %s found for the url - %s", remote, r.remoteUrl), global.StatusInfo)
		return true
	}
	return false
}

func NewGetRemoteName(repo *git2go.Repository, remoteUrl string, validate Validation) Name {
	return remoteName{
		repo:      repo,
		remoteUrl: remoteUrl,
		validate:  validate,
	}
}
