package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
)

// ResetAllItems removes all the staged items back to the staging area.
//
// go-git fails to reset all changes for a newly initialized repo, so the function falls back
// to the git client for resetting all items if the go-git Reset fails
func ResetAllItems(repo *git.Repository) string {
	w, wErr := repo.Worktree()

	if wErr != nil {
		logger.Log(fmt.Sprintf("Error occurred while fetching repo worktree -> %s", wErr.Error()), global.StatusError)
		return "STAGE_ALL_REMOVE_FAILED"
	} else {
		logger.Log("Proceeding with removing all changes", global.StatusInfo)
		err := w.Reset(&git.ResetOptions{
			Mode: git.MixedReset,
		})

		if err != nil {
			logger.Log(fmt.Sprintf("Failed while performing git reset -> %s", err.Error()), global.StatusError)
			logger.Log(fmt.Sprintf("Falling back to git client"), global.StatusWarning)
			args := []string{"reset"}
			cmd := utils.GetGitClient(w.Filesystem.Root(), args)
			removeMsg, err := cmd.Output()

			if err != nil {
				return global.RemoveAllItemsError
			} else {
				logger.Log("All staged items have been removed from staging area -> "+string(removeMsg), global.StatusInfo)
				return "STAGE_ALL_REMOVE_SUCCESS"
			}
		} else {
			logger.Log("All staged items have been removed from staging area", global.StatusInfo)
			return "STAGE_ALL_REMOVE_SUCCESS"
		}
	}

}
