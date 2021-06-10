package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"runtime/debug"
)

type ChangedItemsInterface interface {
	ChangedFiles() *model.GitChangeResults
}

type ChangedItemStruct struct {
	Repo     *git2go.Repository
	RepoPath string
}

func checkChangedFilesError(err error) error {
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}
	return nil
}

// ChangedFiles returns the list of changes from the target
// The function organizes the tracked, untracked and staged files in separate slices and returns the struct *model.GitChangeResults
func (c ChangedItemStruct) ChangedFiles() *model.GitChangeResults {
	logger.Log("Fetching the current status of the repo", global.StatusInfo)
	var changedFileList []*string
	var stagedFileList []*string
	var unTrackedList []*string
	var stagedMap = make(map[string]bool)
	var errStatus error

	logger.Log("Fetching changed files from the repo", global.StatusInfo)
	repo := c.Repo

	defer func() {
		if r := recover(); r != nil {
			logger.Log(fmt.Sprintf("%v", r), global.StatusError)
			logger.Log(string(debug.Stack()), global.StatusWarning)
		}
	}()

	head, _ := repo.Head()
	if head == nil {
		logger.Log("Repo has no HEAD. Treating it as a newly initialized repo", global.StatusWarning)
		statusList, statusListErr := repo.StatusList(&git2go.StatusOptions{
			Show:  git2go.StatusShowIndexAndWorkdir,
			Flags: git2go.StatusOptIncludeUntracked,
		})
		errStatus = checkChangedFilesError(statusListErr)

		n, _ := statusList.EntryCount()
		for i := 0; i < n; i++ {
			entry, entryErr := statusList.ByIndex(i)

			if entryErr == nil {
				diff := entry.IndexToWorkdir
				stagedDiff := entry.HeadToIndex

				if stagedDiff.Status.String() != "Unmodified" {
					fileEntry := stagedDiff.NewFile.Path
					stagedMap[fileEntry] = true
					stagedFileList = append(stagedFileList, &fileEntry)
				}

				if diff.NewFile.Path != "" {
					unTrackedList = append(unTrackedList, &diff.NewFile.Path)
				}
			}
		}
		return &model.GitChangeResults{
			GitUntrackedFiles: unTrackedList,
			GitChangedFiles:   nil,
			GitStagedFiles:    stagedFileList,
		}
	}

	commit, commitErr := repo.LookupCommit(head.Target())
	errStatus = checkChangedFilesError(commitErr)

	tree, treeErr := commit.Tree()
	errStatus = checkChangedFilesError(treeErr)

	diff, diffErr := repo.DiffTreeToWorkdirWithIndex(tree, nil)
	errStatus = checkChangedFilesError(diffErr)

	repoIndex, indexErr := repo.Index()
	errStatus = checkChangedFilesError(indexErr)

	stagedDiff, stagedDiffErr := repo.DiffTreeToIndex(tree, repoIndex, nil)
	errStatus = checkChangedFilesError(stagedDiffErr)

	statusList, statusListErr := repo.StatusList(&git2go.StatusOptions{
		Show:     git2go.StatusShowWorkdirOnly,
		Flags:    git2go.StatusOptIncludeUntracked,
		Pathspec: nil,
	})
	errStatus = checkChangedFilesError(statusListErr)

	stat, statErr := stagedDiff.Stats()
	errStatus = checkChangedFilesError(statErr)

	stagedFileCount := stat.FilesChanged()

	if stagedFileCount > 0 {
		n, _ := stagedDiff.NumDeltas()
		for d := 0; d < n; d++ {
			delta, deltaErr := stagedDiff.Delta(d)
			errStatus = checkChangedFilesError(deltaErr)

			fileEntry := delta.NewFile.Path
			logger.Log(fmt.Sprintf("Staged file --> %v", fileEntry), global.StatusInfo)
			stagedMap[fileEntry] = true
			stagedFileList = append(stagedFileList, &fileEntry)
		}
	} else {
		logger.Log("No staged Files", global.StatusWarning)
		stagedFileList = nil
	}

	var n = 0
	n, _ = diff.NumDeltas()
	for d := 0; d < n; d++ {
		delta, deltaErr := diff.Delta(d)
		errStatus = checkChangedFilesError(deltaErr)

		fileEntry := delta.NewFile.Path
		if !stagedMap[fileEntry] {
			changeStatus := delta.Status
			if changeStatus != git2go.DeltaUntracked {
				logger.Log(fmt.Sprintf("Changed file --> %v", delta.NewFile.Path), global.StatusInfo)
				entry := delta.Status.String()[0:1] + "," + delta.NewFile.Path
				changedFileList = append(changedFileList, &entry)
			}
		} else {
			changeFlag := delta.NewFile.Flags
			if changeFlag == 8 {
				logger.Log(fmt.Sprintf("Staged and Changed file --> %v", delta.NewFile.Path), global.StatusInfo)
				entry := delta.Status.String()[0:1] + "," + delta.NewFile.Path
				changedFileList = append(changedFileList, &entry)
			}
		}
	}

	numStatus, numStatusErr := statusList.EntryCount()
	errStatus = checkChangedFilesError(numStatusErr)

	for i := 0; i < numStatus; i++ {
		statusEntry, _ := statusList.ByIndex(i)
		delta := statusEntry.IndexToWorkdir
		if delta.Status == git2go.DeltaUntracked {
			logger.Log(fmt.Sprintf("Untracked file --> %v", delta.NewFile.Path), global.StatusInfo)
			entry := delta.NewFile.Path
			unTrackedList = append(unTrackedList, &entry)
		}
	}

	if errStatus == nil {
		return &model.GitChangeResults{
			GitUntrackedFiles: unTrackedList,
			GitChangedFiles:   changedFileList,
			GitStagedFiles:    stagedFileList,
		}
	} else {
		return &model.GitChangeResults{
			GitUntrackedFiles: nil,
			GitChangedFiles:   nil,
			GitStagedFiles:    nil,
		}
	}
}
