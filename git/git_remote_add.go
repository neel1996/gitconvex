package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/neel1996/gitconvex-server/global"
)

// AddRemote adds a new remote to the target git repo
func AddRemote(repo *git.Repository, remoteName string, remoteURL string) string {
	logger := global.Logger{}
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
