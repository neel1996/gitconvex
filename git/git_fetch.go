package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

type FetchInterface interface {
	FetchFromRemote() *model.FetchResult
}

type FetchStruct struct {
	Repo         *git2go.Repository
	RemoteName   string
	RepoPath     string
	RemoteURL    string
	RemoteBranch string
	RepoName     string
	AuthOption   string
	UserName     string
	Password     string
	SSHKeyPath   string
}

// FetchFromRemote performs a git fetch for the supplied remote and branch (e.g. `git fetch origin main`)
// If the remoteBranch is empty, then a fetch is performed with no branch name (similar to `git fetch`)
func (f FetchStruct) FetchFromRemote() *model.FetchResult {
	repo := f.Repo
	remoteBranch := f.RemoteBranch
	remoteName := f.RemoteName

	localRefSpec := "+refs/heads/" + remoteBranch
	targetRefPsec := "refs/remotes/" + remoteName + "/" + remoteBranch
	targetRemote, _ := repo.Remotes.Lookup(remoteName)

	if targetRemote == nil {
		logger.Log("Target remote is unavailable", global.StatusError)
		return &model.FetchResult{
			Status:       global.FetchFromRemoteError,
			FetchedItems: nil,
		}
	}

	var remoteCallbackObject RemoteCallbackInterface
	remoteCallbackObject = &RemoteCallbackStruct{
		RepoName:   f.RepoName,
		UserName:   f.UserName,
		Password:   f.Password,
		SSHKeyPath: f.SSHKeyPath,
		AuthOption: f.AuthOption,
	}

	fetchOption := &git2go.FetchOptions{
		RemoteCallbacks: remoteCallbackObject.RemoteCallbackSelector(),
	}
	logger.Log(fmt.Sprintf("Fetching changes from -> %s - %s", remoteName, localRefSpec+":"+targetRefPsec), global.StatusInfo)
	err := targetRemote.Fetch([]string{localRefSpec + ":" + targetRefPsec}, fetchOption, "")

	remoteRef, remoteRefErr := repo.References.Lookup(targetRefPsec)
	if remoteRefErr == nil {
		remoteCommit, _ := repo.AnnotatedCommitFromRef(remoteRef)
		if remoteCommit != nil {
			mergeAnalysis, _, mergeErr := repo.MergeAnalysis([]*git2go.AnnotatedCommit{remoteCommit})
			if mergeErr != nil {
				logger.Log("Fetch failed - "+mergeErr.Error(), global.StatusError)
				return &model.FetchResult{
					Status:       global.FetchFromRemoteError,
					FetchedItems: nil,
				}
			} else {
				if mergeAnalysis&git2go.MergeAnalysisUpToDate != 0 {
					logger.Log("No new changes to fetch from remote", global.StatusWarning)
					msg := "No new changes to fetch from remote"
					return &model.FetchResult{
						Status:       global.FetchNoNewChanges,
						FetchedItems: []*string{&msg},
					}
				}
			}
		}

		if err != nil {
			logger.Log("Fetch failed - "+err.Error(), global.StatusError)
			return &model.FetchResult{
				Status:       global.FetchFromRemoteError,
				FetchedItems: nil,
			}
		} else {
			msg := "Changes fetched from remote " + remoteName
			logger.Log(msg, global.StatusInfo)
			return &model.FetchResult{
				Status:       global.FetchFromRemoteSuccess,
				FetchedItems: []*string{&msg},
			}
		}
	} else {
		logger.Log("Fetch failed - "+remoteRefErr.Error(), global.StatusError)
		return &model.FetchResult{
			Status:       global.FetchFromRemoteError,
			FetchedItems: nil,
		}
	}
}
