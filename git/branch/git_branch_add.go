package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type AddBranch interface {
	AddBranch() string
}

type addBranch struct {
	Repo         *git2go.Repository
	BranchName   string
	RemoteSwitch bool
	TargetCommit *git2go.Commit
}

// AddBranch adds a new branch to the target repo
func (a addBranch) AddBranch() string {
	targetCommit := a.TargetCommit
	repo := a.Repo
	branchName := a.BranchName
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

func NewAddBranch(repo *git2go.Repository, branchName string, remoteSwitch bool, targetCommit *git2go.Commit) AddBranch {
	return addBranch{
		Repo:         repo,
		BranchName:   branchName,
		RemoteSwitch: remoteSwitch,
		TargetCommit: targetCommit,
	}
}
