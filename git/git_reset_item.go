package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
)

type ResetInterface interface {
	RemoveItem() string
}

type ResetStruct struct {
	Repo     *git2go.Repository
	RepoPath string
	FileItem string
}

func removeErr(fileItem string, errMsg string) string {
	logger.Log(fmt.Sprintf("Error occurred while removeing item %s -> %s", fileItem, errMsg), global.StatusError)
	return global.RemoveItemError
}

// RemoveItem performs a git rest 'file' to remove the item from the staged area
func (r ResetStruct) RemoveItem() string {
	fileItem := r.FileItem
	repo := r.Repo

	head, headErr := repo.Head()
	if headErr != nil {
		logger.Log("Repo has no HEAD", global.StatusWarning)
		idx, idxErr := repo.Index()

		if idxErr != nil {
			return removeErr(fileItem, idxErr.Error())
		}

		if idxRemoveErr := idx.RemoveByPath(fileItem); idxRemoveErr != nil {
			return removeErr(fileItem, idxRemoveErr.Error())
		}

		if writeErr := idx.Write(); writeErr != nil {
			return removeErr(fileItem, writeErr.Error())
		}
		return global.RemoveItemSuccess
	}

	commit, commitErr := repo.LookupCommit(head.Target())
	if commitErr != nil {
		return removeErr(fileItem, commitErr.Error())
	}

	resetErr := repo.ResetDefaultToCommit(commit, []string{fileItem})
	if resetErr != nil {
		return removeErr(fileItem, resetErr.Error())
	} else {
		return global.RemoveItemSuccess
	}
}
