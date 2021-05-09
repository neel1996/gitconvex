package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"go/types"
)

type CloneInterface interface {
	CloneRepo() (*model.ResponseModel, error)
}

type CloneStruct struct {
	RepoName   string
	RepoPath   string
	RepoURL    string
	AuthOption string
	UserName   string
	Password   string
	SSHKeyPath string
}

// CloneRepo clones the remote repo to the target directory
// It supports options for SSH and HTTPS authentications
func (c CloneStruct) CloneRepo() (*model.ResponseModel, error) {
	authOption := c.AuthOption
	repoPath := c.RepoPath
	repoURL := c.RepoURL
	userName := c.UserName
	password := c.Password
	sshKeyPath := c.SSHKeyPath

	logger.Log(fmt.Sprintf("Initiating repo clone with - %v auth option", authOption), global.StatusInfo)
	var err error
	var r *git2go.Repository

	var remoteCBObject RemoteCallbackInterface
	remoteCBObject = &RemoteCallbackStruct{
		RepoName:   c.RepoName,
		UserName:   userName,
		Password:   password,
		AuthOption: authOption,
		SSHKeyPath: sshKeyPath,
	}
	var remoteCallbacks git2go.RemoteCallbacks
	remoteCallbacks = remoteCBObject.RemoteCallbackSelector()

	r, err = git2go.Clone(repoURL, repoPath, &git2go.CloneOptions{
		FetchOptions: &git2go.FetchOptions{
			RemoteCallbacks: remoteCallbacks,
		},
	})

	if err != nil {
		logger.Log(fmt.Sprintf("Error occurred while cloning repo \n%v", err), global.StatusError)
		return nil, types.Error{Msg: "Git repo clone failed"}
	}

	logger.Log(fmt.Sprintf("Repo %v - Cloned to target directory - %s", c.RepoName, r.Workdir()), global.StatusInfo)

	return &model.ResponseModel{
		Status:    "success",
		Message:   "Git clone completed",
		HasFailed: false,
	}, nil
}
