package git

import (
	"fmt"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
)

type ResetInterface interface {
	RemoveItem() string
}

type ResetStruct struct {
	RepoPath string
	FileItem string
}

func removeErr(fileItem string, errMsg string) string {
	logger := global.Logger{}
	logger.Log(fmt.Sprintf("Error occurred while removeing item %s -> %s", fileItem, errMsg), global.StatusError)
	return global.RemoveItemError
}

// RemoveItem performs a git rest 'file' to remove the item from the staged area
// Uses the gitclient module, as go-git does not support selective reset
func (r ResetStruct) RemoveItem() string {
	repoPath := r.RepoPath
	fileItem := r.FileItem
	args := []string{"reset", fileItem}
	cmd := utils.GetGitClient(repoPath, args)

	if cmd.String() == "" {
		return removeErr(fileItem, "Error occurred while fetching git client")
	} else {
		removeMsg, err := cmd.Output()

		if err != nil {
			return removeErr(fileItem, err.Error())
		} else {
			logger.Log(fmt.Sprintf("Staged item -> %s removed from the staging area\n%s", fileItem, string(removeMsg)), global.StatusInfo)
			return global.RemoveItemSuccess
		}
	}
}
