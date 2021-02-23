package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
)

type PushInterface interface {
	PushToRemote() string
}

type PushStruct struct {
	Repo         *git2go.Repository
	RepoName     string
	AuthOption   string
	UserName     string
	Password     string
	SSHKeyPath   string
	RemoteName   string
	RemoteBranch string
	RepoPath     string
}

// PushToRemote pushed the commits to the selected remote repository
// By default it will choose the current branch and pushes to the matching remote branch
func (p PushStruct) PushToRemote() string {
	repo := p.Repo
	remoteBranch := p.RemoteBranch
	remoteName := p.RemoteName
	targetRefPsec := "refs/heads/" + remoteBranch

	remote, remoteErr := repo.Remotes.Lookup(remoteName)
	if remoteErr != nil {
		logger.Log(remoteErr.Error(), global.StatusError)
		return global.PushToRemoteError
	}

	var remoteCallbackObject RemoteCallbackInterface
	remoteCallbackObject = &RemoteCallbackStruct{
		RepoName:   p.RepoName,
		UserName:   p.UserName,
		Password:   p.Password,
		SSHKeyPath: p.SSHKeyPath,
		AuthOption: p.AuthOption,
	}

	pushOption := &git2go.PushOptions{
		RemoteCallbacks: remoteCallbackObject.RemoteCallbackSelector(),
	}
	err := remote.Push([]string{targetRefPsec}, pushOption)

	if err != nil {
		logger.Log(fmt.Sprintf("Error occurred while pushing changes to -> %s : %s\n%s", remoteName, targetRefPsec, err.Error()), global.StatusError)
		return global.PushToRemoteError
	} else {
		logger.Log(fmt.Sprintf("commits pushed to remote -> %s : %s", remoteName, targetRefPsec), global.StatusInfo)
		return global.PushToRemoteSuccess
	}
}
