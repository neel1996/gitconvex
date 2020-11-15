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
	"go/types"
	"io"
)

// CloneHandler clones the remote repo to the target directory
// It supports options for SSH and HTTPS authentications

func CloneHandler(repoPath string, repoURL string, authOption string, userName *string, password *string) (*model.ResponseModel, error) {
	logger := global.Logger{}

	gitSSHAuth, authErr := ssh.NewSSHAgentAuth("git")

	if authOption == "ssh" && authErr != nil {
		logger.Log(authErr.Error(), global.StatusError)
		return nil, authErr
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
		logger.Log(fmt.Sprintf("Reop %v - Cloned to target directory", r), global.StatusInfo)

		return &model.ResponseModel{
			Status:    "success",
			Message:   "Git clone completed",
			HasFailed: false,
		}, nil
	}
}
