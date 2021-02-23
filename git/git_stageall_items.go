package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
)

type StageAllInterface interface {
	StageAllItems() string
}

type StageAllStruct struct {
	Repo *git.Repository
}

// StageAllItems stages all the modified and untracked items to the worktree
//
// The function relies on the native git client to stage an item, as go-git staging is time consuming for huge repos
func (s StageAllStruct) StageAllItems() string {
	repo := s.Repo
	w, wErr := repo.Worktree()
	logger.Log(fmt.Sprintf("Staging all changes from repo"), global.StatusInfo)

	if wErr != nil {
		logger.Log(fmt.Sprintf("Error occurred while fetching worktree -> %s", wErr.Error()), global.StatusError)
		return global.StageAllItemsError
	} else {
		logger.Log("Proceeding with staging all changes", global.StatusInfo)
		args := []string{"add", "--all"}
		cmd := utils.GetGitClient(w.Filesystem.Root(), args)
		cmdStr, cmdErr := cmd.Output()
		if cmdErr != nil {
			logger.Log(fmt.Sprintf("Error occurred while fetching repo status -> %s", cmdErr.Error()), global.StatusError)
			return global.StageAllItemsError
		} else {
			logger.Log(fmt.Sprintf("All changes staged -> %s", cmdStr), global.StatusInfo)
			return global.StageAllSuccess
		}
	}
}
