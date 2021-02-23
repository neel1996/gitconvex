package git

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
	"io"
	"strings"
)

type PushInterface interface {
	PushToRemote() string
	windowsPush() string
}

type PushStruct struct {
	Repo         *git.Repository
	RemoteName   string
	RemoteBranch string
	RepoPath     string
}

// windowsPush is used for pushing changes using the git client if the platform is windows
// go-git push fails in windows due to SSH authentication error
func (p PushStruct) windowsPush() string {
	repoPath := p.RepoPath
	remoteName := p.RemoteName
	branch := p.RemoteBranch

	args := []string{"push", "-u", remoteName, branch}
	cmd := utils.GetGitClient(repoPath, args)
	cmdStr, cmdErr := cmd.Output()

	if cmdErr != nil {
		logger.Log(fmt.Sprintf("Push failed -> %s", cmdErr.Error()), global.StatusError)
		return global.PushToRemoteError
	} else {
		logger.Log(fmt.Sprintf("Changes pushed to remote -> %s", cmdStr), global.StatusInfo)
		return global.PushToRemoteSuccess
	}
}

// PushToRemote pushed the commits to the selected remote repository
// By default it will choose the current branch and pushes to the matching remote branch
func (p PushStruct) PushToRemote() string {
	repo := p.Repo
	remoteBranch := p.RemoteBranch
	remoteName := p.RemoteName

	targetRefPsec := "refs/heads/" + remoteBranch + ":refs/heads/" + remoteBranch
	w, _ := repo.Worktree()

	b := new(bytes.Buffer)
	sshAuth, sshErr := ssh.NewSSHAgentAuth("git")
	logger.Log(fmt.Sprintf("Pushing changes to remote -> %s : %s", remoteName, targetRefPsec), global.StatusInfo)

	if sshErr != nil {
		logger.Log(fmt.Sprintf("Authentication failed -> %s", sshErr.Error()), global.StatusError)

		if w == nil {
			return global.PushToRemoteError
		}
		logger.Log("Falling back to native git client for pushing changes", global.StatusWarning)
		return p.windowsPush()
	}

	remote, remoteErr := repo.Remote(remoteName)
	if remoteErr != nil {
		logger.Log(remoteErr.Error(), global.StatusError)
		return global.PushToRemoteError
	}

	err := remote.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   []config.RefSpec{config.RefSpec(targetRefPsec)},
		Auth:       sshAuth,
		Progress: sideband.Progress(func(f io.Writer) io.Writer {
			return f
		}(b)),
	})

	if err != nil {
		if strings.Contains(err.Error(), "ssh: handshake failed: ssh:") {
			logger.Log("push failed. Retrying push with git client", global.StatusWarning)
			return p.windowsPush()
		}
		logger.Log(fmt.Sprintf("Error occurred while pushing changes to -> %s : %s\n%s", remoteName, targetRefPsec, err.Error()), global.StatusError)
		return global.PushToRemoteError
	} else {
		logger.Log(fmt.Sprintf("commits pushed to remote -> %s : %s\n%v", remoteName, targetRefPsec, b.String()), global.StatusInfo)
		return global.PushToRemoteSuccess
	}
}
