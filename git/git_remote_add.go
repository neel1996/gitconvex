package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type AddRemoteInterface interface {
	AddRemote() *model.RemoteMutationResult
}

type AddRemoteStruct struct {
	Repo       *git2go.Repository
	RemoteName string
	RemoteURL  string
}

// AddRemote adds a new remote to the target git repo
func (a AddRemoteStruct) AddRemote() *model.RemoteMutationResult {
	repo := a.Repo
	remoteName := a.RemoteName
	remoteURL := a.RemoteURL

	remote, err := repo.Remotes.Create(remoteName, remoteURL)

	if err == nil {
		logger.Log(fmt.Sprintf("New remote %s added to the repo", remote.Name()), global.StatusInfo)
		return &model.RemoteMutationResult{Status: global.RemoteAddSuccess}
	} else {
		logger.Log("Remote addition Failed -> "+err.Error(), global.StatusError)
		return &model.RemoteMutationResult{Status: global.RemoteAddError}
	}
}
