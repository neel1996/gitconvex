package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
)

func ResetAllItems(repo *git.Repository) string {
	w, wErr := repo.Worktree()

	if wErr != nil {
		logger.Log(fmt.Sprintf("Error occurred while fetching repo worktree -> %s", wErr.Error()), global.StatusError)
		return "STAGE_ALL_REMOVE_FAILED"
	} else {
		err := w.Reset(&git.ResetOptions{
			Mode: git.MixedReset,
		})

		if err != nil {
			logger.Log(fmt.Sprintf("Failed while performing git reset -> %s", err.Error()), global.StatusError)
			return "STAGE_ALL_REMOVE_FAILED"
		} else {
			logger.Log("All staged items have been removed from staging area", global.StatusInfo)
			return "STAGE_ALL_REMOVE_SUCCESS"
		}
	}

}
