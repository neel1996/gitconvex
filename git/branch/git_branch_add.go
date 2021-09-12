package branch

//go:generate git/branch/git_branch_add.go

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/validator"
)

type Add interface {
	AddBranch(branchName string, remoteSwitch bool, targetCommit *git2go.Commit) error
}

type addBranch struct {
	repo             middleware.Repository
	branchName       string
	remoteSwitch     bool
	targetCommit     *git2go.Commit
	branchValidation validator.ValidatorWithStringFields
}

func (a addBranch) AddBranch(branchName string, remoteSwitch bool, targetCommit *git2go.Commit) error {
	a.branchName = branchName
	a.remoteSwitch = remoteSwitch
	a.targetCommit = targetCommit

	err := a.validateAddBranchFields()
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}

	head, headErr := a.repo.Head()
	if headErr != nil {
		logger.Log(fmt.Sprintf("Unable to fetch HEAD -> %s", headErr.Error()), global.StatusError)
		return headErr
	}

	targetCommit, validationErr := a.validateTargetCommit(a.targetCommit, a.repo, head)
	if validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return validationErr
	}

	logger.Log(fmt.Sprintf("Adding new branch -> %s", a.branchName), global.StatusInfo)

	_, branchErr := a.repo.CreateBranch(a.branchName, targetCommit, false)
	if branchErr != nil {
		logger.Log(fmt.Sprintf("Failed to add branch - %s - %s", a.branchName, branchErr.Error()), global.StatusError)
		return branchErr
	}

	logger.Log(fmt.Sprintf("Added new branch - %s to the repo", a.branchName), global.StatusInfo)
	return nil
}

func (a addBranch) validateAddBranchFields() error {
	err := a.branchValidation.ValidateWithFields(a.branchName)
	if err != nil {
		return err
	}
	return nil
}

func (a addBranch) validateTargetCommit(targetCommit *git2go.Commit, repo middleware.Repository, head middleware.Reference) (*git2go.Commit, error) {
	if targetCommit != nil {
		return targetCommit, nil
	}

	headCommit, headCommitErr := repo.LookupCommit(head.Target())
	if headCommit == nil {
		return nil, headCommitErr
	}

	return headCommit, nil
}

func NewAddBranch(repo middleware.Repository, branchValidation validator.ValidatorWithStringFields) Add {
	return addBranch{
		repo:             repo,
		branchValidation: branchValidation,
	}
}
