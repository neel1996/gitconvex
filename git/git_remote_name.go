package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
)

// GetRemoteName function returns the name of the remote based on the supplied remote URL

func GetRemoteName(repo *git.Repository, remoteURL string) string {
	var remoteName string
	logger := global.Logger{}

	remotes, remoteErr := repo.Remotes()

	if remoteErr != nil {
		logger.Log(remoteErr.Error(), global.StatusError)
	} else {
		for _, remote := range remotes {
			if remote.Config().URLs[0] == remoteURL {
				remoteName = remote.Config().Name
			}
		}
	}
	return remoteName
}
