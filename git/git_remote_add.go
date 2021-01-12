package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/neel1996/gitconvex-server/global"
)

type AddRemoteInterface interface {
	AddRemote() string
}

type AddRemoteStruct struct {
	Repo       *git.Repository
	RemoteName string
	RemoteURL  string
}

// AddRemote adds a new remote to the target git repo
func (a AddRemoteStruct) AddRemote() string {
	logger := global.Logger{}

	repo := a.Repo
	remoteName := a.RemoteName
	remoteURL := a.RemoteURL

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{remoteURL},
	})

	if err == nil {
		logger.Log("Remoted addition completed for"+remote.String(), global.StatusInfo)
		return global.RemoteAddSuccess
	} else {
		logger.Log("Remote addition Failed!", global.StatusError)
		return global.RemoteAddError
	}
}
