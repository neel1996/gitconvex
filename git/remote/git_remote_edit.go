package remote

import (
	"errors"
	"fmt"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
)

type Edit interface {
	EditRemote() error
}

type editRemote struct {
	repo             middleware.Repository
	remoteName       string
	remoteURL        string
	remoteValidation Validation
	remoteList       List
}

func (e editRemote) EditRemote() error {
	repo := e.repo

	validationErr := e.remoteValidation.ValidateRemoteFields(e.remoteName, e.remoteURL)
	if validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return validationErr
	}

	_, listErr := repo.Remotes().List()
	if listErr != nil {
		logger.Log(listErr.Error(), global.StatusError)
		return listErr
	}

	if remoteAvailabilityErr := e.isRemotePresentInRepo(); remoteAvailabilityErr != nil {
		logger.Log(remoteAvailabilityErr.Error(), global.StatusError)
		return remoteAvailabilityErr
	}

	err := repo.Remotes().SetUrl(e.remoteName, e.remoteURL)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}

	logger.Log("Remote data has been updated successfully", global.StatusInfo)
	return nil
}

func (e editRemote) isRemotePresentInRepo() error {
	var err error
	remoteList := e.remoteList.GetAllRemotes()

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

func NewEditRemote(repo middleware.Repository, remoteName string, remoteURL string, remoteValidation Validation, remoteList List) Edit {
	return editRemote{
		repo:             repo,
		remoteName:       remoteName,
		remoteURL:        remoteURL,
		remoteValidation: remoteValidation,
		remoteList:       remoteList,
	}
}
