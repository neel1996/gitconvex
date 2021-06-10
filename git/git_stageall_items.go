package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"strings"
)

type StageAllInterface interface {
	StageAllItems() string
}

type StageAllStruct struct {
	Repo *git2go.Repository
}

// StageAllItems stages all the modified and untracked items to the worktree
//
// The function relies on the native git client to stage an item, as go-git staging is time consuming for huge repos
func (s StageAllStruct) StageAllItems() string {
	var filesForStaging []*string
	var stageAllSwitch = true

	repo := s.Repo

	var stageItemObject StageItemInterface

	var changedFiesObject ChangedItemsInterface
	changedFiesObject = ChangedItemStruct{
		Repo: repo,
	}

	changedFiles := changedFiesObject.ChangedFiles()
	modifiedFiles := changedFiles.GitChangedFiles
	unTrackedFiles := changedFiles.GitUntrackedFiles

	if changedFiles != nil {
		if len(modifiedFiles) > 0 {
			for _, modifiedEntry := range modifiedFiles {
				splitEntry := strings.Split(*modifiedEntry, ",")
				if len(splitEntry) > 0 {
					filesForStaging = append(filesForStaging, &splitEntry[1])
				}
			}
		}

		if len(unTrackedFiles) > 0 {
			filesForStaging = append(filesForStaging, unTrackedFiles...)
		}

		for _, stageItem := range filesForStaging {
			stageItemObject = StageItemStruct{
				Repo:     repo,
				FileItem: *stageItem,
			}
			logger.Log("From stage all : Staging item -> "+*stageItem, global.StatusInfo)
			stageStatus := stageItemObject.StageItem()
			if stageStatus != global.StageItemSuccess {
				logger.Log("Stage individual item failed => "+*stageItem, global.StatusError)
				stageAllSwitch = false
			}
		}
	} else {
		return global.StageAllItemsError
	}

	if stageAllSwitch {
		logger.Log(fmt.Sprintf("All changes have been staged"), global.StatusInfo)
		return global.StageAllSuccess
	} else {
		logger.Log(fmt.Sprintf("Staging all items failed"), global.StatusError)
		return global.StageAllItemsError
	}
}
