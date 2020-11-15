package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
)

func addError(fileItem string, errMsg string) string {
	logger := global.Logger{}
	logger.Log(fmt.Sprintf("Error occurred while staging %s -> %s", fileItem, errMsg), global.StatusError)
	return "ADD_ITEM_FAILED"
}

// StageItem stages a selected file from the target repo
// The function relies on the native git client to stage an item, as go-git staging is time consuming for huge repos
func StageItem(repo *git.Repository, fileItem string) string {
	logger := global.Logger{}
	w, wErr := repo.Worktree()
	if wErr != nil {
		return addError(fileItem, wErr.Error())
	} else {
		args := []string{"add", fileItem}
		cmd := utils.GetGitClient(w.Filesystem.Root(), args)
		cmdString, cmdErr := cmd.Output()
		if cmdErr != nil {
			logger.Log(fmt.Sprintf("Staging of %s failed -> %s", fileItem, cmdErr.Error()), global.StatusError)
			return addError(fileItem, cmdErr.Error())
		} else {
			logger.Log(fmt.Sprintf("File -> %s staged \n%s", fileItem, cmdString), global.StatusInfo)
			return "ADD_ITEM_SUCCESS"
		}
	}
}
