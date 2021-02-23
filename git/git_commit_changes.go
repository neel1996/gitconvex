package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"runtime/debug"
	"strings"
)

type CommitInterface interface {
	CommitChanges() string
}

type CommitStruct struct {
	Repo          *git2go.Repository
	CommitMessage string
	RepoPath      string
}

func checkCommitError(err error) bool {
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		panic(err)
		return true
	}
	return false
}

// CommitChanges commits the staged changes to the repo
// Rewrites the repo index tree with the staged files to commit the changes
func (c CommitStruct) CommitChanges() string {
	var errStatus bool
	commitMessage := c.CommitMessage
	repo := c.Repo

	defer func() string {
		if r := recover(); r != nil {
			logger.Log(string(debug.Stack()), global.StatusError)
			return global.CommitChangeError
		}
		return ""
	}()

	var formattedMessage = commitMessage
	signature, signatureErr := repo.DefaultSignature()

	if signatureErr != nil {
		logger.Log(fmt.Sprintf("Error occurred while fetching repo signature -> %s", signatureErr.Error()), global.StatusError)
		logger.Log("Setup you user name and email using `git config user.name and git config user.email`", global.StatusWarning)
		return global.CommitChangeError
	}

	if strings.Contains(commitMessage, "||") {
		splitMessage := strings.Split(commitMessage, "||")
		formattedMessage = strings.Join(splitMessage, "\n")
	}

	head, headErr := repo.Head()
	errStatus = checkCommitError(headErr)

	headCommit, headCommitErr := repo.LookupCommit(head.Target())
	errStatus = checkCommitError(headCommitErr)

	repoIndex, indexErr := repo.Index()
	errStatus = checkCommitError(indexErr)

	newTreeId, writeTreeErr := repoIndex.WriteTree()
	errStatus = checkCommitError(writeTreeErr)

	newTree, newTreeErr := repo.LookupTree(newTreeId)
	errStatus = checkCommitError(newTreeErr)

	newCommitId, err := repo.CreateCommit("HEAD", signature, signature, formattedMessage, newTree, headCommit)
	errStatus = checkCommitError(err)

	_, newCommitErr := repo.LookupCommit(newCommitId)
	errStatus = checkCommitError(newCommitErr)

	head, headErr = repo.Head()
	errStatus = checkCommitError(headErr)

	newRef, newRefErr := head.SetTarget(newCommitId, formattedMessage)
	errStatus = checkCommitError(newRefErr)

	if errStatus {
		return global.CommitChangeError
	} else {
		logger.Log(fmt.Sprintf("New commit %s created", newCommitId.String()), global.StatusInfo)
		logger.Log(newRef.Name(), global.StatusInfo)
		return global.CommitChangeSuccess
	}
}
