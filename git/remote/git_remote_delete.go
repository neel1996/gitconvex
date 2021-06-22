package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

type Delete interface {
	DeleteRemote() *model.RemoteMutationResult
}

type deleteRemote struct {
	repo       *git2go.Repository
	remoteName string
}

func (d deleteRemote) DeleteRemote() *model.RemoteMutationResult {
	remoteToBeDeleted := d.remoteName

	deletionError := &model.RemoteMutationResult{Status: global.RemoteDeleteError}

	err := d.deleteSelectedRemote(d.remoteName)
	if err != nil {
		logger.Log(fmt.Sprintf("Remote => %s cannot be found in the repo", remoteToBeDeleted), global.StatusError)
		return deletionError
	}

	return &model.RemoteMutationResult{Status: global.RemoteDeleteSuccess}
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

func NewDeleteRemoteInterface(repo *git2go.Repository, remoteName string) Delete {
	return deleteRemote{
		repo:       repo,
		remoteName: remoteName,
	}
}
