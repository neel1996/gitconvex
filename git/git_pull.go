package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type PullInterface interface {
	PullFromRemote() *model.PullResult
}

type PullStruct struct {
	Repo         *git2go.Repository
	RemoteURL    string
	RemoteBranch string
	RemoteName   string
	RepoName     string
	RepoPath     string
	AuthOption   string
	UserName     string
	Password     string
	SSHKeyPath   string
}

func returnPullErr(msg string) *model.PullResult {
	logger.Log(msg, global.StatusError)
	return &model.PullResult{
		Status:      global.PullFromRemoteError,
		PulledItems: nil,
	}
}

// PullFromRemote pulls the changes from the remote repository using the remote URL and branch name received
func (p PullStruct) PullFromRemote() *model.PullResult {
	repo := p.Repo
	remoteBranch := p.RemoteBranch
	remoteURL := p.RemoteURL

	var remoteDataObject RemoteDataInterface
	remoteDataObject = RemoteDataStruct{
		Repo:      repo,
		RemoteURL: remoteURL,
	}

	remoteName := remoteDataObject.GetRemoteName()
	localRefSpec := "+refs/heads/" + remoteBranch
	targetRefPsec := "refs/remotes/" + remoteName + "/" + remoteBranch
	targetRemote, _ := repo.Remotes.Lookup(remoteName)

	if targetRemote == nil {
		return returnPullErr("Target remote is unavailable")
	}

	var remoteCallbackObject RemoteCallbackInterface
	remoteCallbackObject = &RemoteCallbackStruct{
		RepoName:   p.RepoName,
		UserName:   p.UserName,
		Password:   p.Password,
		SSHKeyPath: p.SSHKeyPath,
		AuthOption: p.AuthOption,
	}
	fetchOption := &git2go.FetchOptions{
		RemoteCallbacks: remoteCallbackObject.RemoteCallbackSelector(),
		UpdateFetchhead: true,
	}

	logger.Log(fmt.Sprintf("Fetching changes from -> %s - %s", remoteName, targetRefPsec), global.StatusInfo)
	fetchErr := targetRemote.Fetch([]string{localRefSpec + ":" + targetRefPsec}, fetchOption, "")
	if fetchErr != nil {
		return returnPullErr("Fetch Failed : " + fetchErr.Error())
	}

	remoteRef, remoteRefErr := repo.References.Lookup(targetRefPsec)
	if remoteRefErr == nil {
		remoteCommit, _ := repo.LookupCommit(remoteRef.Target())
		fmt.Println(remoteRef.Name())
		fmt.Println(remoteCommit.Id(), remoteCommit.Message())

		annotatedCommit, _ := repo.AnnotatedCommitFromRef(remoteRef)
		if annotatedCommit != nil {
			mergeAnalysis, _, mergeErr := repo.MergeAnalysis([]*git2go.AnnotatedCommit{annotatedCommit})
			if mergeErr != nil {
				return returnPullErr("Pull failed - " + mergeErr.Error())
			} else {
				if mergeAnalysis&git2go.MergeAnalysisUpToDate != 0 {
					logger.Log("No new changes to pull from remote", global.StatusWarning)
					msg := "No new changes to pull from remote"
					return &model.PullResult{
						Status:      global.PullNoNewChanges,
						PulledItems: []*string{&msg},
					}
				} else {
					err := repo.Merge([]*git2go.AnnotatedCommit{annotatedCommit}, nil, &git2go.CheckoutOptions{
						Strategy: git2go.CheckoutSafe,
					})
					if err != nil {
						return returnPullErr("Annotated merge failed : " + err.Error())
					} else {
						repoIndex, _ := repo.Index()
						if repoIndex.HasConflicts() {
							return returnPullErr("Conflicts encountered while pulling changes")
						}

						indexTree, indexTreeErr := repoIndex.WriteTree()
						if indexTreeErr != nil {
							return returnPullErr("Index Tree Error : " + indexTreeErr.Error())
						}
						remoteTree, treeErr := repo.LookupTree(indexTree)
						if treeErr != nil {
							return returnPullErr("Tree Error : " + treeErr.Error())
						}
						checkoutErr := repo.CheckoutTree(remoteTree, &git2go.CheckoutOptions{
							Strategy: git2go.CheckoutForce,
						})
						if checkoutErr != nil {
							return returnPullErr("Tree checkout error : " + checkoutErr.Error())
						}

						localRef, localRefErr := repo.LookupBranch(remoteBranch, git2go.BranchLocal)
						if localRefErr != nil {
							return returnPullErr("Local Reference lookup error :" + localRefErr.Error())
						}
						_, targetErr := localRef.SetTarget(remoteRef.Target(), "")
						if targetErr != nil {
							return returnPullErr("Target set error : " + targetErr.Error())
						}
						head, _ := repo.Head()
						if head == nil {
							return returnPullErr("HEAD is nil")
						}

						headTarget, _ := head.SetTarget(remoteCommit.Id(), "")
						if headTarget == nil {
							return returnPullErr("Unable to set target to HEAD")
						}

						logger.Log("New changes pulled from remote -> "+targetRemote.Name(), global.StatusInfo)
						msg := "New changes pulled from remote"
						return &model.PullResult{
							Status:      global.PullFromRemoteSuccess,
							PulledItems: []*string{&msg},
						}
					}
				}
			}
		}
	} else {
		return returnPullErr(remoteRefErr.Error())
	}
	return returnPullErr("")
}
