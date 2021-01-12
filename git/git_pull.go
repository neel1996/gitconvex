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
	"github.com/neel1996/gitconvex-server/utils"
	"go/types"
	"io"
	"strings"
)

type PullInterface interface {
	PullFromRemote() *model.PullResult
	windowsPull() *model.PullResult
}

type PullStruct struct {
	Repo         *git.Repository
	RemoteURL    string
	RemoteBranch string
	RemoteName   string
	RepoPath     string
}

// windowsPull is used for pulling changes using the git client if the platform is windows
// go-git pull fails in windows due to SSH authentication error
func (p PullStruct) windowsPull() *model.PullResult {
	remoteName := p.RemoteName
	repoPath := p.RepoPath
	branch := p.RemoteBranch

	args := []string{"pull", remoteName, branch}
	cmd := utils.GetGitClient(repoPath, args)
	cmdStr, cmdErr := cmd.Output()

	if cmdErr != nil {
		logger.Log(fmt.Sprintf("Pull failed -> %s", cmdErr.Error()), global.StatusError)

		return &model.PullResult{
			Status:      global.PullFromRemoteError,
			PulledItems: nil,
		}
	} else {
		if strings.Contains(string(cmdStr), "Already up to date") {
			logger.Log(fmt.Sprintf("No new changes available -> %s", cmdStr), global.StatusInfo)

			msg := "No changes to pull from " + remoteName
			return &model.PullResult{
				Status:      global.PullNoNewChanges,
				PulledItems: []*string{&msg},
			}
		}
		msg := "New Items Pulled from remote " + remoteName
		logger.Log(fmt.Sprintf("Changes pulled from remote -> %s", cmdStr), global.StatusInfo)
		return &model.PullResult{
			Status:      global.PullFromRemoteSuccess,
			PulledItems: []*string{&msg},
		}
	}
}

// PullFromRemote pulls the changes from the remote repository using the remote URL and branch name received
func (p PullStruct) PullFromRemote() *model.PullResult {
	var pullErr error
	logger := global.Logger{}

	repo := p.Repo
	remoteURL := p.RemoteURL
	remoteBranch := p.RemoteBranch

	var remoteDataObject RemoteDataInterface
	remoteDataObject = RemoteDataStruct{
		Repo:      repo,
		RemoteURL: remoteURL,
	}

	remoteName := remoteDataObject.GetRemoteName()
	p.RemoteName = remoteName

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

		if remoteName == "" {
			return &model.PullResult{
				Status:      global.PullFromRemoteError,
				PulledItems: nil,
			}
		}

		if sshErr != nil {
			logger.Log("Authentication method failed -> "+sshErr.Error(), global.StatusError)
			w, _ := repo.Worktree()
			if w == nil {
				return &model.PullResult{
					Status:      global.PullFromRemoteError,
					PulledItems: nil,
				}
			}
			return p.windowsPull()
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
			logger.Log("Pull failed with error : "+pullErr.Error(), global.StatusWarning)
			msg := "No changes to pull from " + remoteName
			return &model.PullResult{
				Status:      global.PullNoNewChanges,
				PulledItems: []*string{&msg},
			}
		} else {
			if strings.Contains(pullErr.Error(), "ssh: handshake failed: ssh:") || strings.Contains(pullErr.Error(), "invalid auth method") {
				logger.Log("Pull failed. Retrying pull with git client", global.StatusWarning)
				return p.windowsPull()
			}
			logger.Log(pullErr.Error(), global.StatusError)
			return &model.PullResult{
				Status:      global.PullFromRemoteError,
				PulledItems: nil,
			}
		}
	} else {
		logger.Log("New items pulled from remote", global.StatusInfo)
		msg := "New Items Pulled from remote " + remoteName
		return &model.PullResult{
			Status:      global.PullFromRemoteSuccess,
			PulledItems: []*string{&msg},
		}
	}
}
