package commit

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"strings"
)

type Changes interface {
	Add() error
}

type changes struct {
	repo          middleware.Repository
	commitMessage []string
}

// Add commits the staged changes to the repo
// Rewrites the repo index tree with the staged files to commit the changes
func (c changes) Add() error {
	var (
		headCommit       *git2go.Commit
		formattedMessage string
		repo             = c.repo
	)

	signature, err := c.validateSignature(repo)
	if err != nil {
		return err
	}

	formattedMessage = c.formatCommitMessage()

	headCommit, headCommitErr := c.currentHeadCommit()
	if headCommitErr != nil {
		logger.Log(headCommitErr.Error(), global.StatusError)
		return headCommitErr
	}

	newTree, treeErr := c.treeForNewCommit(repo)
	if treeErr != nil {
		logger.Log(treeErr.Error(), global.StatusError)
		return treeErr
	}

	newCommitId, newCommitIdErr := c.createNewCommit(headCommit, repo, signature, formattedMessage, newTree)
	if newCommitIdErr != nil {
		logger.Log(newCommitIdErr.Error(), global.StatusError)
		return newCommitIdErr
	}

	setHeadErr := c.setHeadToNewCommit(repo, newCommitId, formattedMessage)
	if setHeadErr != nil {
		logger.Log(setHeadErr.Error(), global.StatusError)
		return setHeadErr
	}

	logger.Log(fmt.Sprintf("New commit %s created", newCommitId.String()), global.StatusInfo)
	return nil
}

func (c changes) validateSignature(repo middleware.Repository) (*git2go.Signature, error) {
	signature, signatureErr := repo.DefaultSignature()
	if signatureErr != nil {
		logger.Log(fmt.Sprintf("Error occurred while fetching repo signature -> %s", signatureErr.Error()), global.StatusError)
		logger.Log("Setup you user name and email using `git config user.name and git config user.email`", global.StatusWarning)
		return nil, signatureErr
	}

	return signature, nil
}

func (c changes) setHeadToNewCommit(repo middleware.Repository, newCommitId *git2go.Oid, formattedMessage string) error {
	if _, newCommitErr := repo.LookupCommit(newCommitId); newCommitErr != nil {
		logger.Log(newCommitErr.Error(), global.StatusError)
		return newCommitErr
	}

	head, headErr := repo.Head()
	if headErr != nil {
		logger.Log(headErr.Error(), global.StatusError)
		return headErr
	}

	_, newRefErr := head.SetTarget(newCommitId, formattedMessage)
	if newRefErr != nil {
		logger.Log(newRefErr.Error(), global.StatusError)
		return newRefErr
	}

	return nil
}

func (c changes) createNewCommit(headCommit *git2go.Commit, repo middleware.Repository, signature *git2go.Signature, formattedMessage string, newTree *git2go.Tree) (*git2go.Oid, error) {
	var (
		newCommitId    *git2go.Oid
		newCommitIdErr error
	)

	if headCommit == nil {
		newCommitId, newCommitIdErr = repo.CreateCommit("HEAD", signature, signature, formattedMessage, newTree)
	} else {
		newCommitId, newCommitIdErr = repo.CreateCommit("HEAD", signature, signature, formattedMessage, newTree, headCommit)
	}

	return newCommitId, newCommitIdErr
}

func (c changes) treeForNewCommit(repo middleware.Repository) (*git2go.Tree, error) {
	treeId, indexErr := c.updateRepoIndex()
	if indexErr != nil {
		logger.Log(indexErr.Error(), global.StatusError)
		return nil, indexErr
	}

	newTree, newTreeErr := repo.LookupTree(treeId)
	if newTreeErr != nil {
		logger.Log(newTreeErr.Error(), global.StatusError)
		return nil, newTreeErr
	}
	return newTree, nil
}

func (c changes) updateRepoIndex() (*git2go.Oid, error) {
	repoIndex, indexErr := c.repo.Index()
	if indexErr != nil {
		return nil, indexErr
	}

	return repoIndex.WriteTree()
}

func (c changes) currentHeadCommit() (*git2go.Commit, error) {
	head, headErr := c.repo.Head()
	if headErr != nil {
		logger.Log("Repo has no HEAD. Proceeding with a NIL HEAD", global.StatusWarning)
		return nil, nil
	}

	return c.repo.LookupCommit(head.Target())
}

func (c changes) formatCommitMessage() string {
	if len(c.commitMessage) > 0 {
		return strings.Join(c.commitMessage, "\n")
	}

	return strings.Join(c.commitMessage, "")
}

func NewCommitChanges(repo middleware.Repository, message []string) Changes {
	return changes{
		repo:          repo,
		commitMessage: message,
	}
}
