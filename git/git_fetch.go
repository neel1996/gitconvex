package git

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"io"
)

// FetchFromRemote performs a git fetch for the supplied remote and branch (e.g. `git fetch origin main`)
// If the remoteBranch is empty, then a fetch is performed with no branch name (similar to `git fetch`)

func FetchFromRemote(repo *git.Repository, remoteURL string, remoteBranch string) *model.FetchResult {
	remoteName := GetRemoteName(repo, remoteURL)
	logger := global.Logger{}

	targetRefPsec := "refs/heads/" + remoteBranch + ":refs/remotes/" + remoteBranch
	b := new(bytes.Buffer)
	var fetchErr error
	gitSSHAuth, sshErr := ssh.NewSSHAgentAuth("git")

	if sshErr != nil {
		logger.Log("Authentication method failed -> "+sshErr.Error(), global.StatusError)
		return &model.FetchResult{
			Status:       "FETCH ERROR",
			FetchedItems: nil,
		}
	}

	logger.Log(fmt.Sprintf("Fetching changes from -> %s : %s", remoteURL, targetRefPsec), global.StatusInfo)

	if remoteURL != "" && remoteBranch != "" {
		fetchErr = repo.Fetch(&git.FetchOptions{
			RemoteName: remoteName,
			Auth:       gitSSHAuth,
			RefSpecs:   []config.RefSpec{config.RefSpec(targetRefPsec)},
			Progress: sideband.Progress(func(f io.Writer) io.Writer {
				return f
			}(b)),
		})
	} else {
		fetchErr = repo.Fetch(&git.FetchOptions{
			RemoteName: git.DefaultRemoteName,
			Auth:       gitSSHAuth,
			Progress: sideband.Progress(func(f io.Writer) io.Writer {
				return f
			}(b)),
		})
	}

	if fetchErr != nil {
		if fetchErr.Error() == "already up-to-date" {
			logger.Log(fetchErr.Error(), global.StatusWarning)
			return &model.FetchResult{
				Status:       "NEW CHANGES ABSENT",
				FetchedItems: nil,
			}
		} else {
			logger.Log(fetchErr.Error(), global.StatusError)
			return &model.FetchResult{
				Status:       "FETCH ERROR",
				FetchedItems: nil,
			}
		}

	} else {
		logger.Log(b.String(), global.StatusInfo)
		logger.Log("Changes fetched from remote", global.StatusInfo)
		msg := fmt.Sprintf("Changes fetched from %v", git.DefaultRemoteName)
		return &model.FetchResult{
			Status:       "CHANGES FETCHED FROM REMOTE",
			FetchedItems: []*string{&msg},
		}
	}

}
