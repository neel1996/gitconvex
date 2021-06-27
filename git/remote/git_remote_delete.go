package remote

import (
	"errors"
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type Delete interface {
	DeleteRemote() error
}

type deleteRemote struct {
	repo       *git2go.Repository
	remoteName string
	validate   Validation
}

// DeleteRemote deletes the remote based on the specified remoteName
func (d deleteRemote) DeleteRemote() error {
	if validationError := d.validateRemoteFields(); validationError != nil {
		return validationError
	}

	err := d.deleteSelectedRemote(d.remoteName)
	if err != nil {
		logger.Log(fmt.Sprintf("Remote => %s cannot be found in the repo", d.remoteName), global.StatusError)
		return err
	}

	return nil
}

func (d deleteRemote) validateRemoteFields() error {
	if d.repo == nil {
		return errors.New("repo is nil")
	}

	if d.remoteName == "" {
		return errors.New("remote name cannot be empty")
	}
	return nil
}

func (d *deleteRemote) deleteSelectedRemote(remoteEntry string) error {
	err := d.repo.Remotes.Delete(remoteEntry)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}

	logger.Log(fmt.Sprintf("The remote => %s deleted from repo", d.remoteName), global.StatusInfo)
	return nil
}

func NewDeleteRemote(repo *git2go.Repository, remoteName string, validate Validation) Delete {
	return deleteRemote{
		repo:       repo,
		remoteName: remoteName,
		validate:   validate,
	}
}
