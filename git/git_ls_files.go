package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/neel1996/gitconvex-server/utils"
	"go/types"
	"io/ioutil"
	"strings"
	"sync"
)

// This go file relies on git installed on the host or the git client packed with the build application -> ./gitclient{.exe}
// Git client dependency was induced as the go-git based log traversal was highly time consuming

type LsFileInfo struct {
	Content           []*string
	Commits           []*string
	TotalTrackedCount *int
}

type dirCommitDataModel struct {
	dirNameList   []*string
	dirCommitList []*string
}

type fileCommitDataModel struct {
	fileNameList   []*string
	fileCommitList []*string
}

var logger global.Logger
var selectedDir string
var waitGroup sync.WaitGroup

// pathFilterCheck validates if the path held by the log iterator is tracked by the repo

func pathFilterCheck(filterPath string) bool {
	if strings.Contains(filterPath, selectedDir) {
		return true
	}
	return false
}

// dirCommitHandler collects the commit messages for the directories present in the target repo

func DirCommitHandler(dirName *string, repoPath string, fileChan chan string, commitChan chan string, waitGroup *sync.WaitGroup) {
	args := []string{"log", "--oneline", "-1", "--pretty=format:%s", *dirName}
	cmd := utils.GetGitClient(repoPath, args)

	if cmd.String() != "" {
		dirStr := *dirName + ":directory"

		commitLog, err := cmd.Output()
		if err != nil {
			logger.Log(fmt.Sprintf("Command execution for -> {{%s}} failed with error %v", cmd.String(), err.Error()), global.StatusError)
			fmt.Println(commitLog)
		} else {
			commitMsg := string(commitLog)
			logger.Log(fmt.Sprintf("Fetching commits for file -> %s --> %s", *dirName, commitLog), global.StatusInfo)
			fileChan <- dirStr
			commitChan <- commitMsg
			waitGroup.Done()
		}
	}
}

// fileCommitHandler collects the commit messages for the files present in the target repo

func FileCommitHandler(file *string, repoPath string, fileChan chan string, commitChan chan string, waitGroup *sync.WaitGroup) {
	args := []string{"log", "--oneline", "-1", "--pretty=format:%s", *file}
	cmd := utils.GetGitClient(repoPath, args)

	if cmd.String() != "" {
		var fileStr string

		if strings.Contains(*file, "/") {
			splitEntry := strings.Split(*file, "/")
			fileStr = splitEntry[len(splitEntry)-1] + ":File"
		} else {
			fileStr = *file + ":File"
		}

		commitLog, err := cmd.Output()
		if err != nil {
			logger.Log(err.Error(), global.StatusError)
		} else {
			commitMsg := string(commitLog)
			logger.Log(fmt.Sprintf("Fetching commits for file -> %v --> %s", *file, commitLog), global.StatusInfo)

			fileChan <- fileStr
			commitChan <- commitMsg
			waitGroup.Done()
		}
	}
}

// TrackedFileCount returns the total number of files tracked by the target git repo

func TrackedFileCount(repo *git.Repository, trackedFileCountChan chan int) {
	var totalFileCount int
	logger := global.Logger{}

	head, _ := repo.Head()
	hash := head.Hash()

	allCommits, _ := repo.CommitObject(hash)
	tObj, _ := allCommits.Tree()

	err := tObj.Files().ForEach(func(file *object.File) error {
		if file != nil {
			totalFileCount++
			return nil
		} else {
			return types.Error{Msg: "File from the tree is empty"}
		}
	})
	tObj.Files().Close()

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		trackedFileCountChan <- 0
	} else {
		logger.Log(fmt.Sprintf("Total Tracked Files : %v", totalFileCount), global.StatusInfo)
		trackedFileCountChan <- totalFileCount
	}
	close(trackedFileCountChan)
}

// ListFiles collects the list of tracked files and their latest respective commit messages
// Used to visualize the git repo in the front-end file explorer in a github explorer based fashion

func ListFiles(repo *git.Repository, repoPath string, directoryName string) *model.GitFolderContentResults {
	logger := global.Logger{}
	logger.Log("Collecting tracked file list from the repo", global.StatusInfo)

	var targetPath string
	var fileList []*string
	var dirList []*string
	var fileFilterList []*string
	var commitList []*string

	fileFilterList = nil
	commitList = nil

	if directoryName != "" {
		targetPath = repoPath + "/" + directoryName
	} else {
		targetPath = repoPath
	}

	content, _ := ioutil.ReadDir(targetPath)

	for _, files := range content {
		var fileName string
		if directoryName != "" {
			fileName = directoryName + "/" + files.Name()
		} else {
			fileName = files.Name()
		}
		if files.IsDir() && fileName != ".git" {
			dirName := fileName
			dirList = append(dirList, &dirName)
		} else {
			if fileName != ".git" {
				fileList = append(fileList, &fileName)
			}
		}
	}
	content = nil

	var fileListChan = make(chan string)
	var commitListChan = make(chan string)

	for _, file := range fileList {
		waitGroup.Add(1)
		go FileCommitHandler(file, repoPath, fileListChan, commitListChan, &waitGroup)
		fileName := <-fileListChan
		commitMsg := <-commitListChan

		if commitMsg != "" {
			fileFilterList = append(fileFilterList, &fileName)
			commitList = append(commitList, &commitMsg)
		}
	}

	for _, dir := range dirList {
		waitGroup.Add(1)
		go DirCommitHandler(dir, repoPath, fileListChan, commitListChan, &waitGroup)
		fileName := <-fileListChan
		commitMsg := <-commitListChan

		if commitMsg != "" {
			fileFilterList = append(fileFilterList, &fileName)
			commitList = append(commitList, &commitMsg)
		}
	}

	waitGroup.Wait()
	close(fileListChan)
	close(commitListChan)

	if len(fileFilterList) == 0 || len(commitList) == 0 {
		msg := "NO_TRACKED_FILES"
		noFileList := []*string{&msg}
		return &model.GitFolderContentResults{
			TrackedFiles:     noFileList,
			FileBasedCommits: commitList,
		}
	}

	return &model.GitFolderContentResults{
		TrackedFiles:     fileFilterList,
		FileBasedCommits: commitList,
	}
}
