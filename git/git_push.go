package git

import (
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/neel1996/gitconvex-server/global"
	"io"
)

func PushToRemote(repo *git.Repository, remoteName string, remoteBranch string) string {
	targetRefPsec := "refs/heads/" + remoteBranch + ":refs/heads/" + remoteBranch
	b := new(bytes.Buffer)
	sshAuth, sshErr := ssh.NewSSHAgentAuth("git")
	logger.Log(fmt.Sprintf("Pushing changes to remote -> %s : %s", remoteName, targetRefPsec), global.StatusInfo)

	if sshErr != nil {
		logger.Log(fmt.Sprintf("Authentication failed -> %s", sshErr.Error()), global.StatusError)
		return "PUSH_FAILED"
	}

	remote, remoteErr := repo.Remote(remoteName)
	if remoteErr != nil {
		logger.Log(remoteErr.Error(), global.StatusError)
		return "PUSH_FAILED"
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
		logger.Log(fmt.Sprintf("Error occurred while pushing changes to -> %s : %s\n%s", remoteName, targetRefPsec, err.Error()), global.StatusError)
		return "PUSH_FAILED"
	} else {
		logger.Log(fmt.Sprintf("commits pushed to remote -> %s : %s\n%v", remoteName, targetRefPsec, b.String()), global.StatusInfo)
		return "PUSH_SUCCESS"
	}
}
