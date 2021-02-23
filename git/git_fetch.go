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
	"github.com/neel1996/gitconvex-server/utils"
	"io"
)

type FetchInterface interface {
	FetchFromRemote() *model.FetchResult
	windowsFetch() *model.FetchResult
}

type FetchStruct struct {
	Repo         *git.Repository
	RemoteName   string
	RepoPath     string
	RemoteURL    string
	RemoteBranch string
}

// windowsFetch is used for fetching changes using the git client if the platform is windows
// go-git fetch fails in windows due to SSH authentication error
func (f FetchStruct) windowsFetch() *model.FetchResult {
	var args []string

	remoteName := f.RemoteName
	repoPath := f.RepoPath
	branch := f.RemoteBranch

	if remoteName == "" && branch == "" {
		args = []string{"fetch"}
	} else {
		branchReference := branch + ":" + branch
		args = []string{"fetch", remoteName, branchReference}
	}
	cmd := utils.GetGitClient(repoPath, args)
	cmdStr, cmdErr := cmd.Output()

	if cmdErr != nil {
		logger.Log(fmt.Sprintf("Fetch failed -> %s", cmdErr.Error()), global.StatusError)

		return &model.FetchResult{
			Status:       global.FetchFromRemoteError,
			FetchedItems: nil,
		}
	} else {
		logger.Log(fmt.Sprintf("Changes fetched from remote - %s\n%s", remoteName, cmdStr), global.StatusInfo)

		msg := fmt.Sprintf("Changes fetched from remote %v", remoteName)
		return &model.FetchResult{
			Status:       "CHANGES FETCHED FROM REMOTE",
			FetchedItems: []*string{&msg},
		}
	}
}

// FetchFromRemote performs a git fetch for the supplied remote and branch (e.g. `git fetch origin main`)
// If the remoteBranch is empty, then a fetch is performed with no branch name (similar to `git fetch`)
func (f FetchStruct) FetchFromRemote() *model.FetchResult {
	repo := f.Repo
	remoteURL := f.RemoteURL
	remoteBranch := f.RemoteBranch
	repoPath := f.RepoPath

	var remoteDataObject RemoteDataInterface
	remoteDataObject = RemoteDataStruct{
		Repo:      repo,
		RemoteURL: remoteURL,
	}

	remoteName := remoteDataObject.GetRemoteName()
	logger := global.Logger{}

	targetRefPsec := "refs/heads/" + remoteBranch + ":refs/remotes/" + remoteBranch
	b := new(bytes.Buffer)
	var fetchErr error
	gitSSHAuth, sshErr := ssh.NewSSHAgentAuth("git")
	w, _ := repo.Worktree()

	// Check if repo path is empty and fetch path from worktree
	if repoPath == "" {
		repoPath = w.Filesystem.Root()
	}

	if sshErr != nil {
		logger.Log("Authentication method failed -> "+sshErr.Error(), global.StatusError)
		if w == nil {
			return &model.FetchResult{
				Status:       global.FetchFromRemoteError,
				FetchedItems: nil,
			}
		}
		logger.Log("Retrying fetch with fallback module using git client", global.StatusWarning)
		return f.windowsFetch()
	}

	logger.Log(fmt.Sprintf("Fetching changes from -> %s : %s", remoteURL, targetRefPsec), global.StatusInfo)

	if remoteURL != "" && remoteBranch != "" {
		if remoteName == "" {
			return &model.FetchResult{
				Status:       global.FetchFromRemoteError,
				FetchedItems: nil,
			}
		}

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
				Status:       global.FetchNoNewChanges,
				FetchedItems: nil,
			}
		} else {
			logger.Log(fetchErr.Error(), global.StatusError)
			logger.Log("Fetch failed. Retrying fetch with git client", global.StatusWarning)
			return f.windowsFetch()
		}

	} else {
		logger.Log(b.String(), global.StatusInfo)
		logger.Log("Changes fetched from remote", global.StatusInfo)

		msg := fmt.Sprintf("Changes fetched from remote %v", remoteName)
		return &model.FetchResult{
			Status:       global.FetchFromRemoteSuccess,
			FetchedItems: []*string{&msg},
		}
	}

}
