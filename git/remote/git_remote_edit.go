package remote

import (
	"errors"
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
	validate   Validation
}

func (e editRemote) EditRemote() error {
	repo := e.repo

	validationErr := e.validateRemoteFields()
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

	err := repo.Remotes.SetUrl(e.remoteName, e.remoteURL)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}

	logger.Log("Remote data has been updated successfully", global.StatusInfo)
	return nil
}

func (e editRemote) validateRemoteFields() error {
	validationError := e.validate.ValidateRemoteFields(e.repo)
	if validationError != nil {
		return validationError
	}

	if e.remoteName == "" || e.remoteURL == "" {
		return errors.New("required field(s) are empty")
	}
	return nil
}

func NewEditRemote(repo *git2go.Repository, remoteName string, remoteURL string, validation Validation) Edit {
	return editRemote{
		repo:       repo,
		remoteName: remoteName,
		remoteURL:  remoteURL,
		validate:   validation,
	}
}
