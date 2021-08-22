package remote

import (
	"fmt"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
)

type Add interface {
	NewRemote() error
}

type addRemote struct {
	repo             middleware.Repository
	remoteName       string
	remoteURL        string
	remoteValidation Validation
}

// NewRemote adds a new remote to the target git repo
func (a addRemote) NewRemote() error {
	if validationErr := a.remoteValidation.ValidateRemoteFields(a.remoteName, a.remoteURL); validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return validationErr
	}

	remote, err := a.repo.Remotes().Create(a.remoteName, a.remoteURL)
	if err != nil {
		logger.Log("Remote addition Failed -> "+err.Error(), global.StatusError)
		return err
	}

	logger.Log(fmt.Sprintf("New remote %s added to the repo", remote.Name()), global.StatusInfo)
	return nil
}

func NewAddRemote(repo middleware.Repository, remoteName string, remoteURL string, remoteValidation Validation) Add {
	return addRemote{
		repo:             repo,
		remoteName:       remoteName,
		remoteURL:        remoteURL,
		remoteValidation: remoteValidation,
	}
}
