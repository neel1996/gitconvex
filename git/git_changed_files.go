package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/neel1996/gitconvex-server/utils"
	"strings"
)

/*
	The file invokes the GetGitClient function to run native git commands in the system terminal to fetch some results.
	The (WorkTree) Status function is slower for huge repos which is taking minutes to return the results, so git client is called
	to do the job
*/

type ChangeInterface interface {
	GetUntrackedFiles(untrackedChan chan []*string)
	GetModifiedFiles(modifiedFileChan chan []*string)
	GetStagedFiles(stagedFileChan chan []*string)
	ChangedFiles() *model.GitChangeResults
}

type ChangedStruct struct {
	Repo     *git.Repository
	RepoPath string
}

// GetUntrackedFiles executes a native git command to fetch the list of untracked files
func (c ChangedStruct) GetUntrackedFiles(untrackedChan chan []*string) {
	logger.Log("Fetching untracked files from the repo", global.StatusInfo)
	args := []string{"ls-files", "--others", "--exclude-standard"}
	repoPath := c.RepoPath

	cmd := utils.GetGitClient(repoPath, args)
	cmdString, cmdErr := cmd.Output()
	if cmdErr != nil {
		logger.Log(fmt.Sprintf("Command execution failed -> %s", cmdErr.Error()), global.StatusError)
		untrackedChan <- []*string{}
	} else {
		trimStr := strings.TrimSpace(string(cmdString))
		splitLines := strings.Split(trimStr, "\n")
		logger.Log(fmt.Sprintf("Untracked files --> %v", splitLines), global.StatusInfo)
		var fileList []*string
		for _, file := range splitLines {
			untrackedFile := file
			if file != "" {
				fileList = append(fileList, &untrackedFile)
			}
		}
		untrackedChan <- fileList
	}
}

// GetStagedFiles executes a native git command to fetch the list of staged files
func (c ChangedStruct) GetStagedFiles(stagedFileChan chan []*string) {
	logger.Log("Fetching staged files from the repo", global.StatusInfo)
	args := []string{"diff", "--name-only", "--cached"}
	repoPath := c.RepoPath
	cmd := utils.GetGitClient(repoPath, args)
	cmdString, cmdErr := cmd.Output()
	if cmdErr != nil {
		logger.Log(fmt.Sprintf("Command execution failed -> %s", cmdErr.Error()), global.StatusError)
		stagedFileChan <- []*string{}
	} else {
		trimStr := strings.TrimSpace(string(cmdString))
		splitLines := strings.Split(trimStr, "\n")
		logger.Log(fmt.Sprintf("Staged files --> %v", splitLines), global.StatusInfo)

		var fileList []*string
		for _, file := range splitLines {
			if file != "" {
				temp := file
				fileList = append(fileList, &temp)
			}
		}
		stagedFileChan <- fileList
	}
}

// GetModifiedFiles executes a native git command to fetch all the modified files from the repo
func (c ChangedStruct) GetModifiedFiles(modifiedFileChan chan []*string) {
	logger.Log("Fetching changed files from the repo", global.StatusInfo)
	args := []string{"diff", "--raw"}
	repoPath := c.RepoPath
	cmd := utils.GetGitClient(repoPath, args)
	cmdString, cmdErr := cmd.Output()
	if cmdErr != nil && cmdString != nil {
		logger.Log(fmt.Sprintf("Command execution failed -> %s", cmdErr.Error()), global.StatusError)
		modifiedFileChan <- []*string{}
	} else {
		trimStr := strings.TrimSpace(string(cmdString))
		splitLines := strings.Split(trimStr, "\n")
		var fileList []*string
		for _, change := range splitLines {
			if change == "" {
				continue
			}
			fileName := strings.Fields(change)
			logger.Log(fmt.Sprintf("Changed file --> %v", fileName), global.StatusInfo)
			var changeStr string
			if fileName[4] == "D" {
				changeStr = "D," + fileName[5]
			} else {
				changeStr = "M," + fileName[5]
			}
			fileList = append(fileList, &changeStr)
		}
		modifiedFileChan <- fileList
	}
}

// ChangedFiles returns the list of changes from the target
// The function organizes the tracked, untracked and staged files in separate slices and returns the struct *model.GitChangeResults
func (c ChangedStruct) ChangedFiles() *model.GitChangeResults {
	logger := global.Logger{}
	repo := c.Repo

	logger.Log("Fetching the current status of the repo", global.StatusInfo)
	w, _ := repo.Worktree()
	repoPath := w.Filesystem.Root()

	var unTrackedFileChan = make(chan []*string)
	var changedFileChan = make(chan []*string)
	var stagedFileChan = make(chan []*string)
	c.RepoPath = repoPath

	go c.GetModifiedFiles(changedFileChan)
	go c.GetUntrackedFiles(unTrackedFileChan)
	go c.GetStagedFiles(stagedFileChan)

	// Intermediate return value to close the channels and then return the values
	returnVal := &model.GitChangeResults{
		GitUntrackedFiles: <-unTrackedFileChan,
		GitChangedFiles:   <-changedFileChan,
		GitStagedFiles:    <-stagedFileChan,
	}

	defer func() {
		close(stagedFileChan)
		close(changedFileChan)
		close(unTrackedFileChan)
	}()

	return returnVal
}
