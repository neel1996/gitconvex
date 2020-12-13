package git

import (
	"bytes"
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/neel1996/gitconvex-server/utils"
	"go/types"
	"io"
)

// fallbackClone performs a git clone using the native git client
// If the go-git based clone fails due to an authentication issue, then this function will be invoked to perform a clone
func fallbackClone(repoPath string, repoURL string) (*model.ResponseModel, error) {
	args := []string{"clone", repoURL, repoPath}
	cmd := utils.GetGitClient(".", args)
	cmdStr, cmdErr := cmd.Output()

	if cmdErr != nil {
		logger.Log(fmt.Sprintf("Fallback clone failed -> %s", cmdErr.Error()), global.StatusError)
		return nil, cmdErr
	} else {
		logger.Log(fmt.Sprintf("New repo has been cloned to -> %s -> %s", repoPath, cmdStr), global.StatusInfo)
		return &model.ResponseModel{
			Status:    "success",
			Message:   "Git clone completed",
			HasFailed: false,
		}, nil
	}
}

// CloneHandler clones the remote repo to the target directory
// It supports options for SSH and HTTPS authentications
func CloneHandler(repoPath string, repoURL string, authOption string, userName *string, password *string) (*model.ResponseModel, error) {
	logger := global.Logger{}

	gitSSHAuth, authErr := ssh.NewSSHAgentAuth("git")

	if authOption == "ssh" && authErr != nil {
		logger.Log(authErr.Error(), global.StatusError)
		logger.Log("Auth failed. Retrying with native git client based clone", global.StatusWarning)
		return fallbackClone(repoPath, repoURL)
	} else {
		logger.Log(fmt.Sprintf("Initiating repo clone with - %v auth option", authOption), global.StatusInfo)
		var err error
		var r *git.Repository
		b := new(bytes.Buffer)

		if authOption == "ssh" {
			r, err = git.PlainClone(repoPath, false, &git.CloneOptions{
				URL:  repoURL,
				Auth: gitSSHAuth,
				Progress: sideband.Progress(func(f io.Writer) io.Writer {
					return f
				}(b)),
			})
		} else if authOption == "https" {
			r, err = git.PlainClone(repoPath, false, &git.CloneOptions{
				URL: repoURL,
				Auth: &http.BasicAuth{
					Username: *userName,
					Password: *password,
				},
				Progress: sideband.Progress(func(f io.Writer) io.Writer {
					return f
				}(b)),
			})
		} else {
			r, err = git.PlainClone(repoPath, false, &git.CloneOptions{
				URL: repoURL,
				Progress: sideband.Progress(func(f io.Writer) io.Writer {
					return f
				}(b)),
			})
		}

		if err != nil {
			logger.Log(fmt.Sprintf("Error occurred while cloning repo \n%v", err), global.StatusError)
			return nil, types.Error{Msg: "Git repo clone failed"}
		}

		fmt.Println(b.String())

		var repoRoot string
		w, _ := r.Worktree()
		if w != nil {
			repoRoot = w.Filesystem.Root()
		}

		logger.Log(fmt.Sprintf("Repo %v - Cloned to target directory - %s", r, repoRoot), global.StatusInfo)

		return &model.ResponseModel{
			Status:    "success",
			Message:   "Git clone completed",
			HasFailed: false,
		}, nil
	}
}
