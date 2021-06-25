package remote

import (
	"errors"
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type Add interface {
	NewRemote() error
}

type addRemote struct {
	repo       *git2go.Repository
	remoteName string
	remoteURL  string
}

// NewRemote adds a new remote to the target git repo
func (a addRemote) NewRemote() error {
	if validationErr := a.validateRemoteFields(); validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return validationErr
	}

	remote, err := a.repo.Remotes.Create(a.remoteName, a.remoteURL)
	if err != nil {
		logger.Log("Remote addition Failed -> "+err.Error(), global.StatusError)
		return err
	}

	logger.Log(fmt.Sprintf("New remote %s added to the repo", remote.Name()), global.StatusInfo)
	return nil
}

func (a addRemote) validateRemoteFields() error {
	if a.repo == nil {
		return errors.New("repo is nil")
	}

	if a.remoteName == "" {
		return errors.New("remote name cannot be empty")
	}

	if a.remoteURL == "" {
		return errors.New("remote URL is empty")
	}

	return nil
}

func NewAddRemote(repo *git2go.Repository, remoteName string, remoteURL string) Add {
	return addRemote{
		repo:       repo,
		remoteName: remoteName,
		remoteURL:  remoteURL,
	}
}
