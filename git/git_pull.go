package git

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"go/types"
	"io"
)

// PullFromRemote pulls the changes from the remote repository using the remote URL and branch name received
func PullFromRemote(repo *git.Repository, remoteURL string, remoteBranch string) *model.PullResult {
	var pullErr error
	logger := global.Logger{}
	remoteName := GetRemoteName(repo, remoteURL)
	w, _ := repo.Worktree()
	b := new(bytes.Buffer)

	refName := fmt.Sprintf("refs/heads/%s", remoteBranch)

	ref, refErr := repo.Storer.Reference(plumbing.ReferenceName(refName))

	if refErr != nil {
		fmt.Println(refErr.Error())
		pullErr = types.Error{Msg: "branch reference does not exist"}
	} else {
		logger.Log(fmt.Sprintf("Pulling changes from -> %s : %s", remoteURL, ref.Name()), global.StatusInfo)
		gitSSHAuth, sshErr := ssh.NewSSHAgentAuth("git")

		if sshErr != nil {
			logger.Log("Authentication method failed -> "+sshErr.Error(), global.StatusError)
			return &model.PullResult{
				Status:      "PULL ERROR",
				PulledItems: nil,
			}
		}

		pullErr = w.Pull(&git.PullOptions{
			RemoteName:    remoteName,
			Auth:          gitSSHAuth,
			ReferenceName: ref.Name(),
			Progress: sideband.Progress(func(f io.Writer) io.Writer {
				return f
			}(b)),
			SingleBranch: true,
		})
	}

	// Logging the pull message stream sent from the remote server
	logger.Log(b.String(), global.StatusInfo)

	if pullErr != nil {
		if pullErr.Error() == git.NoErrAlreadyUpToDate.Error() {
			logger.Log(pullErr.Error(), global.StatusWarning)
			msg := "No changes to pull from " + remoteName
			return &model.PullResult{
				Status:      "NEW CHANGES ABSENT",
				PulledItems: []*string{&msg},
			}
		} else {
			logger.Log(pullErr.Error(), global.StatusError)
			return &model.PullResult{
				Status:      "PULL ERROR",
				PulledItems: nil,
			}
		}
	} else {
		logger.Log("New items pulled from remote", global.StatusInfo)
		msg := "New Items Pulled from remote " + remoteName
		return &model.PullResult{
			Status:      "PULL SUCCESS",
			PulledItems: []*string{&msg},
		}
	}
}
