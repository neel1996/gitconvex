package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"io/ioutil"
	"strings"
	"sync"
)

type ListFilesInterface interface {
	DirCommitHandler(dir *string)
	FileCommitHandler(file *string)
	TrackedFileCount(trackedFileCountChan chan int)
	ListFiles() *model.GitFolderContentResults
}

type fileDiffStruct struct {
	diffPath string
}

type commitEntry struct {
	fileEntries   []fileDiffStruct
	commitMessage string
}

type ListFilesStruct struct {
	Repo                 *git2go.Repository
	RepoPath             string
	DirectoryName        string
	Commits              []*git2go.Commit
	AllCommitTreeEntries []commitEntry
	FileName             *string
	fileChan             chan string
	commitChan           chan string
	waitGroup            *sync.WaitGroup
}

type LsFileInfo struct {
	Content           []*string
	Commits           []*string
	TotalTrackedCount *int
}

var logger global.Logger
var waitGroup sync.WaitGroup

// DirCommitHandler collects the commit messages for the directories present in the target repo
func (l ListFilesStruct) DirCommitHandler(dirName *string) {
	fileChan := l.fileChan
	commitChan := l.commitChan
	waitGroup := l.waitGroup
	allCommits := l.Commits

	var dirEntry = ""
	var commitMsg = ""

	for _, commit := range allCommits {
		if commit != nil {
			numParents := commit.ParentCount()
			if numParents > 0 {
				parent := commit.Parent(0)
				commitTree, _ := commit.Tree()
				parentTree, _ := parent.Tree()

				if commitTree != nil && parentTree != nil {
					cId, _ := commitTree.EntryByPath(*dirName)
					pId, _ := parentTree.EntryByPath(*dirName)

					if cId != nil && pId != nil {
						if cId.Id.String() != pId.Id.String() {
							dirEntry = *dirName
							commitMsg = commit.Message()
							break
						}
					}

					// This condition is for files that are newly added as part of a commit
					// but not changed in any other commits i.e. Objects with no changes
					if pId == nil && cId != nil {
						dirEntry = *dirName
						commitMsg = commit.Message()
						break
					}
				}
			} else {
				commitTree, _ := commit.Tree()
				if commitTree != nil {
					_, cIdErr := commitTree.EntryByPath(*dirName)
					if cIdErr == nil {
						dirEntry = *dirName
						commitMsg = commit.Message()
						break
					} else {
						logger.Log(cIdErr.Error(), global.StatusError)
					}
				}
			}
		}
	}

	if dirEntry != "" {
		logger.Log(fmt.Sprintf("Fetching commits for directory -> %s --> %s", dirEntry, strings.Split(commitMsg, "\n")[0]), global.StatusInfo)
	}
	dirStr := dirEntry + ":directory"
	fileChan <- dirStr
	commitChan <- commitMsg
	waitGroup.Done()
}

// FileCommitHandler collects the commit messages for the files present in the target repo
func (l ListFilesStruct) FileCommitHandler(file *string) {
	fileChan := l.fileChan
	commitChan := l.commitChan
	waitGroup := l.waitGroup
	allCommits := l.Commits

	var fileStr string
	var commitMsg string
	var fileEntry = ""

	for _, commit := range allCommits {
		if commit != nil {
			numParents := commit.ParentCount()
			if numParents > 0 {
				parent := commit.Parent(0)
				commitTree, _ := commit.Tree()
				parentTree, _ := parent.Tree()

				if commitTree != nil && parentTree != nil {
					cId, _ := commitTree.EntryByPath(*file)
					pId, _ := parentTree.EntryByPath(*file)

					if cId != nil && pId != nil {
						if cId.Id.String() != pId.Id.String() {
							fileEntry = *file
							commitMsg = commit.Message()
							break
						}
					}

					// This condition is for files that are newly added as part of a commit
					// but not changed in any other commits i.e. Objects with no changes
					if pId == nil && cId != nil {
						fileEntry = *file
						commitMsg = commit.Message()
						break
					}
				}
			} else {
				commitTree, _ := commit.Tree()
				if commitTree != nil {
					_, cIdErr := commitTree.EntryByPath(*file)
					if cIdErr == nil {
						fileEntry = *file
						commitMsg = commit.Message()
						break
					} else {
						logger.Log(cIdErr.Error(), global.StatusError)
					}
				} else {
					continue
				}
			}
		}
	}

	if fileEntry != "" {
		logger.Log(fmt.Sprintf("Fetching commits for file -> %v --> %s", *file, strings.Split(commitMsg, "\n")[0]), global.StatusInfo)

		if strings.Contains(*file, "/") {
			splitEntry := strings.Split(*file, "/")
			fileStr = splitEntry[len(splitEntry)-1] + ":File"
		} else {
			fileStr = *file + ":File"
		}
	}

	fileChan <- fileStr
	commitChan <- commitMsg
	waitGroup.Done()
}

// TrackedFileCount returns the total number of files tracked by the target git repo
func (l ListFilesStruct) TrackedFileCount(trackedFileCountChan chan int) {
	var totalFileCount int
	logger := global.Logger{}

	repo := l.Repo

	head, headErr := repo.Head()
	if headErr != nil {
		logger.Log(fmt.Sprintf("Repo head is invalid -> %s", headErr.Error()), global.StatusError)
		trackedFileCountChan <- 0
	} else {
		hash := head.Branch().Target()
		headCommit, _ := repo.LookupCommit(hash)

		if headCommit != nil {
			tree, treeErr := headCommit.Tree()
			if treeErr != nil {
				logger.Log(treeErr.Error(), global.StatusError)
				trackedFileCountChan <- totalFileCount
			} else {
				err := tree.Walk(func(s string, entry *git2go.TreeEntry) int {
					if entry.Id != nil && entry.Type.String() == "Blob" {
						totalFileCount++
						return 1
					}
					return 0
				})

				if err != nil {
					logger.Log(err.Error(), global.StatusError)
					trackedFileCountChan <- 0
				} else {
					logger.Log(fmt.Sprintf("Total Tracked Files : %v", totalFileCount), global.StatusInfo)
					trackedFileCountChan <- totalFileCount
				}
			}
		}
	}
	close(trackedFileCountChan)
}

// ListFiles collects the list of tracked files and their latest respective commit messages
//
// Used to display the git repo structure in the front-end file explorer in a github explorer based fashion
func (l ListFilesStruct) ListFiles() *model.GitFolderContentResults {
	logger := global.Logger{}
	logger.Log("Collecting tracked file list from the repo", global.StatusInfo)

	directoryName := l.DirectoryName
	repoPath := l.RepoPath
	r, _ := git2go.OpenRepository(l.RepoPath)

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

	l.waitGroup = &waitGroup
	l.fileChan = fileListChan
	l.commitChan = commitListChan

	var allCommits []*git2go.Commit
	var commitEntries []commitEntry

	head, headErr := r.Head()

	if headErr == nil {
		commit, _ := r.LookupCommit(head.Target())
		for commit != nil {
			allCommits = append(allCommits, commit)
			commit = commit.Parent(0)
		}
	} else {
		logger.Log(headErr.Error(), global.StatusError)
		return &model.GitFolderContentResults{
			TrackedFiles:     fileFilterList,
			FileBasedCommits: commitList,
		}
	}

	l.Commits = allCommits
	l.AllCommitTreeEntries = commitEntries

	for _, file := range fileList {
		waitGroup.Add(1)
		go l.FileCommitHandler(file)
		fileName := <-fileListChan
		commitMsg := <-commitListChan

		if commitMsg != "" {
			fileFilterList = append(fileFilterList, &fileName)
			commitList = append(commitList, &commitMsg)
		}
	}

	for _, dir := range dirList {
		waitGroup.Add(1)
		go l.DirCommitHandler(dir)
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
