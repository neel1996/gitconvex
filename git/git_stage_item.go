package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
)

type StageItemInterface interface {
	StageItem() string
	addError(errMsg string) string
}

type StageItemStruct struct {
	Repo     *git.Repository
	FileItem string
}

func (s StageItemStruct) addError(errMsg string) string {
	logger := global.Logger{}
	logger.Log(fmt.Sprintf("Error occurred while staging %s -> %s", s.FileItem, errMsg), global.StatusError)
	return global.StageItemError
}

// StageItem stages a selected file from the target repo
// The function relies on the native git client to stage an item, as go-git staging is time consuming for huge repos
func (s StageItemStruct) StageItem() string {
	logger := global.Logger{}

	repo := s.Repo
	fileItem := s.FileItem

	w, wErr := repo.Worktree()
	if wErr != nil {
		return s.addError(wErr.Error())
	} else {
		args := []string{"add", fileItem}
		cmd := utils.GetGitClient(w.Filesystem.Root(), args)
		cmdString, cmdErr := cmd.Output()
		if cmdErr != nil {
			logger.Log(fmt.Sprintf("Staging of %s failed -> %s", fileItem, cmdErr.Error()), global.StatusError)
			return s.addError(cmdErr.Error())
		} else {
			logger.Log(fmt.Sprintf("File -> %s staged \n%s", fileItem, cmdString), global.StatusInfo)
			return global.StageItemSuccess
		}
	}
}
