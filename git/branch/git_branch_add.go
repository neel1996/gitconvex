package branch

import (
	"errors"
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
)

type Add interface {
	AddBranch() error
}

type addBranch struct {
	repo         *git2go.Repository
	branchName   string
	remoteSwitch bool
	targetCommit *git2go.Commit
}

func (a addBranch) AddBranch() error {
	targetCommit := a.targetCommit
	repo := a.repo
	branchName := a.branchName

	err := a.validateAddBranchFields()
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}

	head, headErr := repo.Head()
	if headErr != nil {
		logger.Log(fmt.Sprintf("Unable to fetch HEAD -> %s", headErr.Error()), global.StatusError)
		return headErr
	}

	targetCommit, validationErr := a.validateTargetCommit(targetCommit, repo, head)
	if validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return validationErr
	}

	logger.Log(fmt.Sprintf("Adding new branch -> %s", branchName), global.StatusInfo)

	_, branchErr := repo.CreateBranch(branchName, targetCommit, false)
	if branchErr != nil {
		logger.Log(fmt.Sprintf("Failed to add branch - %s - %s", branchName, branchErr.Error()), global.StatusError)
		return branchErr
	}

	logger.Log(fmt.Sprintf("Added new branch - %s to the repo", branchName), global.StatusInfo)
	return nil
}

func (a addBranch) validateAddBranchFields() error {
	if a.repo == nil {
		return errors.New("repo is nil")
	}

	if a.branchName == "" {
		return errors.New("branch name is empty")
	}
	return nil
}

func (a addBranch) validateTargetCommit(targetCommit *git2go.Commit, repo *git2go.Repository, head *git2go.Reference) (*git2go.Commit, error) {
	if targetCommit != nil {
		return targetCommit, nil
	}

	headCommit, headCommitErr := repo.LookupCommit(head.Target())
	if headCommit == nil {
		return nil, headCommitErr
	}

	return headCommit, nil
}

func NewAddBranch(repo *git2go.Repository, branchName string, remoteSwitch bool, targetCommit *git2go.Commit) Add {
	return addBranch{
		repo:         repo,
		branchName:   branchName,
		remoteSwitch: remoteSwitch,
		targetCommit: targetCommit,
	}
}
