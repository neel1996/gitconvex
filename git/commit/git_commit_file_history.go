package commit

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"reflect"
)

type FileHistory interface {
	Get(commit middleware.Commit) ([]*model.GitCommitFileResult, error)
}

type fileHistory struct {
	repo middleware.Repository
}

func (f fileHistory) Get(commit middleware.Commit) ([]*model.GitCommitFileResult, error) {
	commitHash := commit.Id().String()
	logger.Log(fmt.Sprintf("Fetching file details for commit %v", commitHash), global.StatusInfo)

	prevTree, commitTree, treeErr := f.treesOf(commit)
	if treeErr != nil {
		return fileHistoryError(treeErr)
	}

	return f.diffBetweenTrees(prevTree, commitTree)
}

func (f fileHistory) treesOf(commit middleware.Commit) (*git2go.Tree, *git2go.Tree, error) {
	var treeErr error
	logger.Log("Getting current and previous trees", global.StatusInfo)

	parents := commit.ParentCount()
	if parents == 0 {
		logger.Log("Commit has no parent", global.StatusError)
		return nil, nil, FileHistoryNoParentError
	}

	previousCommit := commit.Parent(0)

	previousTree, treeErr := previousCommit.Tree()
	currentTree, treeErr := commit.Tree()
	if treeErr != nil {
		logger.Log(treeErr.Error(), global.StatusError)
		return nil, nil, FileHistoryTreeError
	}

	return previousTree, currentTree, nil
}

func (f fileHistory) diffBetweenTrees(previousTree *git2go.Tree, currentTree *git2go.Tree) ([]*model.GitCommitFileResult, error) {
	var history []*model.GitCommitFileResult

	logger.Log("Checking diff between trees", global.StatusInfo)
	diff, diffErr := f.repo.DiffTreeToTree(previousTree, currentTree, nil)
	if diffErr != nil {
		return fileHistoryError(diffErr)
	}

	numDelta, numDeltaErr := diff.NumDeltas()
	if numDeltaErr != nil {
		return fileHistoryError(numDeltaErr)
	}

	for d := 0; d < numDelta; d++ {
		delta, _ := diff.Delta(d)
		status := delta.Status.String()
		history = append(history, &model.GitCommitFileResult{
			Type:     status[0:1],
			FileName: delta.NewFile.Path,
		})
	}

	return history, nil
}

func fileHistoryError(err error) ([]*model.GitCommitFileResult, error) {
	logger.Log(err.Error(), global.StatusError)
	if reflect.TypeOf(err) != reflect.TypeOf(Error{}) {
		return nil, FileHistoryError
	}
	return nil, err
}

func NewFileHistory(repo middleware.Repository) FileHistory {
	return fileHistory{repo: repo}
}
