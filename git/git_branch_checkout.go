package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"strings"
)

type BranchCheckoutInterface interface {
	intermediateFetch()
	CheckoutBranch() string
}

type BranchCheckoutInputs struct {
	Repo       *git.Repository
	BranchName string
}

// intermediateFetch performs a remote fetch if the selected checkout branch is a remote branch
func (inputs BranchCheckoutInputs) intermediateFetch() {
	repo := inputs.Repo
	branchName := inputs.BranchName

	logger.Log("Fetching from remote for remote branch -> "+branchName, global.StatusInfo)
	remoteChan := make(chan RemoteDataModel)

	var remoteDataObject RemoteDataInterface
	remoteDataObject = RemoteDataStruct{
		Repo: repo,
	}
	go remoteDataObject.RemoteData(remoteChan)

	remoteData := <-remoteChan
	remoteURL := remoteData.RemoteURL

	remoteDataObject = RemoteDataStruct{
		Repo:      repo,
		RemoteURL: *remoteURL[0],
	}

	var fetchObject FetchInterface
	fetchObject = FetchStruct{
		Repo:         repo,
		RemoteName:   remoteDataObject.GetRemoteName(),
		RepoPath:     "",
		RemoteURL:    *remoteURL[0],
		RemoteBranch: branchName,
	}

	fetchObject.FetchFromRemote()
}

// CheckoutBranch checks out the branchName received as argument
func (inputs BranchCheckoutInputs) CheckoutBranch() string {
	var isRemoteBranch bool
	var referenceBranchName string

	repo := inputs.Repo
	branchName := inputs.BranchName

	logger := global.Logger{}
	w, _ := repo.Worktree()

	if strings.Contains(branchName, "remotes/") {
		splitRef := strings.Split(branchName, "/")
		branchName = splitRef[len(splitRef)-1]
		referenceBranchName = "refs/heads/" + branchName
		isRemoteBranch = true
	} else {
		referenceBranchName = "refs/heads/" + branchName
	}

	// If the branch is a remote branch then a remote fetch will be performed and then the branch checkout will be initiated
	if isRemoteBranch {
		logger.Log(fmt.Sprintf("Branch - %s is a remote branch\nTrying with intermediate remote fetch!", branchName), global.StatusWarning)
		inputs.intermediateFetch()

		checkoutErr := w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(referenceBranchName),
			Force:  true,
		})
		if checkoutErr != nil {
			logger.Log("Checkout failed - "+checkoutErr.Error(), global.StatusError)
			return global.BranchCheckoutError
		}
	}

	checkoutErr := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(referenceBranchName),
		Keep:   true,
	})
	if checkoutErr != nil {
		logger.Log(checkoutErr.Error(), global.StatusError)
		return global.BranchCheckoutError
	}

	logger.Log(fmt.Sprintf("Current branch checked out to -> %s", branchName), global.StatusInfo)
	return fmt.Sprintf("Head checked out to branch - %v", branchName)
}
