package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
)

type ResetAllInterface interface {
	ResetAllItems() string
}

type ResetAllStruct struct {
	Repo *git2go.Repository
}

// ResetAllItems removes all the staged items back to the staging area
func (r ResetAllStruct) ResetAllItems() string {
	repo := r.Repo
	head, headErr := repo.Head()

	if headErr != nil {
		logger.Log(fmt.Sprintf("Reset All operation failed -> %s", headErr.Error()), global.StatusInfo)
		return global.RemoveAllItemsError
	}

	commit, commitErr := repo.LookupCommit(head.Target())
	if commitErr != nil {
		logger.Log(fmt.Sprintf("Reset All operation failed -> %s", commitErr.Error()), global.StatusInfo)
		return global.RemoveAllItemsError
	}

	err := repo.ResetToCommit(commit, git2go.ResetMixed, nil)

	if err != nil {
		logger.Log(fmt.Sprintf("Reset All operation failed -> %s", err.Error()), global.StatusInfo)
		return global.RemoveAllItemsError
	} else {
		return global.ResetAllSuccess
	}
}
