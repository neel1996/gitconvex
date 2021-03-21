package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type RemoteEditInterface interface {
	EditRemoteUrl() *model.RemoteMutationResult
}

type RemoteEditStruct struct {
	Repo          *git2go.Repository
	RemoteName    string
	NewRemoteName string
	RemoteUrl     string
}

func (e RemoteEditStruct) EditRemoteUrl() *model.RemoteMutationResult {
	repo := e.Repo
	remoteName := e.RemoteName
	remoteUrl := e.RemoteUrl

	remoteCollection := repo.Remotes
	_, listErr := remoteCollection.List()

	if listErr != nil {
		logger.Log(listErr.Error(), global.StatusError)
		return &model.RemoteMutationResult{Status: global.RemoteEditError}
	}

	err := repo.Remotes.SetUrl(remoteName, remoteUrl)

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return &model.RemoteMutationResult{Status: global.RemoteEditError}
	} else {
		logger.Log(fmt.Sprintf("New URL has been set for the remote -> %s", remoteName), global.StatusInfo)
		return &model.RemoteMutationResult{Status: global.RemoteEditSuccess}
	}
}
