package remote

import (
	"errors"
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type Edit interface {
	EditRemote() error
}

type editRemote struct {
	repo       *git2go.Repository
	remoteName string
	remoteURL  string
}

func (e editRemote) EditRemote() error {
	repo := e.repo

	validationErr := NewRemoteValidation(e.repo, e.remoteName, e.remoteURL).ValidateRemoteFields()
	if validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return validationErr
	}

	remoteCollection := repo.Remotes
	_, listErr := remoteCollection.List()

	if listErr != nil {
		logger.Log(listErr.Error(), global.StatusError)
		return listErr
	}

	if remoteAvailabilityErr := e.isRemotePresentInRepo(); remoteAvailabilityErr != nil {
		logger.Log(remoteAvailabilityErr.Error(), global.StatusError)
		return remoteAvailabilityErr
	}

	err := repo.Remotes.SetUrl(e.remoteName, e.remoteURL)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}

	logger.Log("Remote data has been updated successfully", global.StatusInfo)
	return nil
}

func (e editRemote) isRemotePresentInRepo() error {
	var err error
	remoteList := NewRemoteList(e.repo).GetAllRemotes()

	if remoteList == nil {
		err = errors.New("no remotes are present in the repo")
		return err
	}

	for _, remote := range remoteList {
		if remote.RemoteName == e.remoteName {
			logger.Log(fmt.Sprintf("Remote - %s is available in the repo", e.remoteName), global.StatusInfo)
			return nil
		}
	}

	err = errors.New(fmt.Sprintf("remote %s is not available in the repo", e.remoteName))
	return err
}

func NewEditRemote(repo *git2go.Repository, remoteName string, remoteURL string) Edit {
	return editRemote{
		repo:       repo,
		remoteName: remoteName,
		remoteURL:  remoteURL,
	}
}
