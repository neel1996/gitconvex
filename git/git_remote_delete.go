package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

type DeleteRemoteInterface interface {
	DeleteRemote() *model.RemoteMutationResult
}

type DeleteRemoteStruct struct {
	Repo       *git2go.Repository
	RemoteName string
}

func (d *DeleteRemoteStruct) DeleteRemote() *model.RemoteMutationResult {
	repo := d.Repo
	remoteToBeDeleted := d.RemoteName

	remoteList, listErr := repo.Remotes.List()
	if listErr != nil {
		logger.Log(listErr.Error(), global.StatusError)
		return &model.RemoteMutationResult{Status: global.RemoteDeleteError}
	}

	for _, remoteEntry := range remoteList {
		if remoteEntry == remoteToBeDeleted {
			err := repo.Remotes.Delete(remoteEntry)
			if err != nil {
				logger.Log(err.Error(), global.StatusError)
				return &model.RemoteMutationResult{Status: global.RemoteDeleteError}
			} else {
				logger.Log(fmt.Sprintf("The remote => %s deleted from repo", remoteToBeDeleted), global.StatusInfo)
				return &model.RemoteMutationResult{Status: global.RemoteDeleteSuccess}
			}
		}
	}

	logger.Log(fmt.Sprintf("Remote => %s cannot be found in the repo", remoteToBeDeleted), global.StatusError)
	return &model.RemoteMutationResult{Status: global.RemoteDeleteError}
}
