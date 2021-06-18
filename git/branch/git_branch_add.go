package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type Add interface {
	AddBranch() string
}

type addBranch struct {
	repo         *git2go.Repository
	branchName   string
	remoteSwitch bool
	targetCommit *git2go.Commit
}

func (a addBranch) AddBranch() string {
	targetCommit := a.targetCommit
	repo := a.repo
	branchName := a.branchName
	head, headErr := repo.Head()

	logger.Log(fmt.Sprintf("Adding new branch -> %s", branchName), global.StatusInfo)
	if headErr != nil {
		logger.Log(fmt.Sprintf("Unable to fetch HEAD -> %s", headErr.Error()), global.StatusError)
		return global.BranchAddError
	}

	targetCommit, hasValidationErr := a.validateTargetCommit(targetCommit, repo, head)
	if hasValidationErr {
		return global.BranchAddError
	}

	_, branchErr := repo.CreateBranch(branchName, targetCommit, false)
	if branchErr != nil {
		logger.Log(fmt.Sprintf("Failed to add branch - %s - %s", branchName, branchErr.Error()), global.StatusError)
		return global.BranchAddError
	}

	logger.Log(fmt.Sprintf("Added new branch - %s to the repo", branchName), global.StatusInfo)
	return global.BranchAddSuccess
}

func (a addBranch) validateTargetCommit(targetCommit *git2go.Commit, repo *git2go.Repository, head *git2go.Reference) (*git2go.Commit, bool) {
	if targetCommit != nil {
		return targetCommit, false
	}

	headCommit, _ := repo.LookupCommit(head.Target())
	if headCommit == nil {
		return nil, true
	}

	return headCommit, false
}

func NewAddBranch(repo *git2go.Repository, branchName string, remoteSwitch bool, targetCommit *git2go.Commit) Add {
	return addBranch{
		repo:         repo,
		branchName:   branchName,
		remoteSwitch: remoteSwitch,
		targetCommit: targetCommit,
	}
}
