package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

type Add interface {
	NewRemote() *model.RemoteMutationResult
}

type addRemote struct {
	repo       *git2go.Repository
	remoteName string
	remoteURL  string
}

// NewRemote adds a new remote to the target git repo
func (a addRemote) NewRemote() *model.RemoteMutationResult {
	repo := a.repo
	remoteName := a.remoteName
	remoteURL := a.remoteURL

	remote, err := repo.Remotes.Create(remoteName, remoteURL)

	if err != nil {
		logger.Log("Remote addition Failed -> "+err.Error(), global.StatusError)
		return &model.RemoteMutationResult{Status: global.RemoteAddError}
	}

	logger.Log(fmt.Sprintf("New remote %s added to the repo", remote.Name()), global.StatusInfo)
	return &model.RemoteMutationResult{Status: global.RemoteAddSuccess}
}

func NewAddRemote(repo *git2go.Repository, remoteName string, remoteURL string) Add {
	return addRemote{
		repo:       repo,
		remoteName: remoteName,
		remoteURL:  remoteURL,
	}
}
