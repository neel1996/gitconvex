package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type StageItemInterface interface {
	StageItem() string
	addError(errMsg string) string
}

type StageItemStruct struct {
	Repo     *git2go.Repository
	FileItem string
}

func (s StageItemStruct) addError(errMsg string) string {
	logger.Log(fmt.Sprintf("Error occurred while staging %s -> %s", s.FileItem, errMsg), global.StatusError)
	return global.StageItemError
}

// StageItem stages a selected file from the target repo
// The function relies on the native git client to stage an item, as go-git staging is time consuming for huge repos
func (s StageItemStruct) StageItem() string {
	repo := s.Repo
	fileItem := s.FileItem
	repoPath := repo.Workdir()

	fileInfo, _ := os.Stat(filepath.Join(repoPath, fileItem))
	if fileInfo != nil && fileInfo.IsDir() {
		itemPath := filepath.Join(repoPath, fileItem)
		logger.Log(fmt.Sprintf("Item %s is a directory", itemPath), global.StatusInfo)
		dirContent, dirReadErr := ioutil.ReadDir(itemPath)
		isStageDirSuccess := true

		if dirReadErr == nil {
			for _, i := range dirContent {
				dirPath := strings.Split(fileItem, repoPath)[0]
				if strings.TrimSpace(dirPath)[len(dirPath)-1:] != "/" {
					dirPath = dirPath + "/"
				}
				s.FileItem = dirPath + i.Name()
				logger.Log(fmt.Sprintf("Staging %s recurssively", s.FileItem), global.StatusInfo)
				if s.StageItem() == global.StageItemError {
					isStageDirSuccess = false
				}
			}

			if isStageDirSuccess {
				return global.StageItemSuccess
			} else {
				return global.StageItemError
			}
		} else {
			logger.Log(fmt.Sprintf("Uanble to read directory directory -> %v", dirReadErr.Error()), global.StatusWarning)
			return s.addError("Empty directory cannot be staged")
		}
	} else {
		repoIndex, repoIndexErr := repo.Index()
		if repoIndexErr != nil {
			return s.addError(repoIndexErr.Error())
		}

		stageErr := repoIndex.AddAll([]string{fileItem}, git2go.IndexAddDefault, nil)
		if stageErr != nil {
			return s.addError(stageErr.Error())
		} else {
			indexWriteErr := repoIndex.Write()
			if indexWriteErr != nil {
				return s.addError(indexWriteErr.Error())
			}

			logger.Log(fmt.Sprintf("File -> %s staged", fileItem), global.StatusInfo)
			return global.StageItemSuccess
		}
	}
}
