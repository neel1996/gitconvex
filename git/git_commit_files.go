package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type CommitFileListInterface interface {
	CommitFileList() []*model.GitCommitFileResult
}

type CommitFileListStruct struct {
	Repo       *git2go.Repository
	CommitHash string
}

// Common function for returning empty response in case of errors
func returnCommitFileError(err error) []*model.GitCommitFileResult {
	logger.Log(err.Error(), global.StatusError)
	return []*model.GitCommitFileResult{{
		Type:     "",
		FileName: "",
	}}
}

// CommitFileList returns the list of files modified, added or removed as part of a commit
// compares the previous commit tree with the current commit tree and, returns the change type (M|D|A) and file name
func (c CommitFileListStruct) CommitFileList() []*model.GitCommitFileResult {
	var res []*model.GitCommitFileResult

	repo := c.Repo
	commitHash := c.CommitHash

	logger.Log(fmt.Sprintf("Fetching file details for commit %v", commitHash), global.StatusInfo)

	oid, oidErr := git2go.NewOid(commitHash)
	if oidErr != nil {
		return returnCommitFileError(oidErr)
	}

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		return returnCommitFileError(err)
	}

	numParents := commit.ParentCount()
	if numParents > 0 {
		prevCommit := commit.Parent(0)
		prevTree, _ := prevCommit.Tree()
		commitTree, _ := commit.Tree()

		if commitTree != nil && prevTree != nil {
			diff, diffErr := repo.DiffTreeToTree(prevTree, commitTree, nil)
			if diffErr != nil {
				return returnCommitFileError(diffErr)
			}

			numDelta, numDeltaErr := diff.NumDeltas()
			if numDeltaErr != nil {
				return returnCommitFileError(numDeltaErr)
			}
			for d := 0; d < numDelta; d++ {
				delta, _ := diff.Delta(d)
				status := delta.Status.String()
				res = append(res, &model.GitCommitFileResult{
					Type:     status[0:1],
					FileName: delta.NewFile.Path,
				})
			}
		}
	} else {
		return []*model.GitCommitFileResult{{
			Type:     "",
			FileName: "",
		}}
	}
	return res
}
