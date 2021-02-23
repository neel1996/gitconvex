package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type CompareCommitInterface interface {
	CompareCommit() []*model.GitCommitFileResult
}

type CompareCommitStruct struct {
	Repo                *git.Repository
	BaseCommitString    string
	CompareCommitString string
}

// CompareCommit compares the tree difference between two commits and returns the file difference status
func (c CompareCommitStruct) CompareCommit() []*model.GitCommitFileResult {
	var res []*model.GitCommitFileResult

	repo := c.Repo
	baseCommitString := c.BaseCommitString
	compareCommitString := c.CompareCommitString

	baseHash := plumbing.NewHash(baseCommitString)
	compareHash := plumbing.NewHash(compareCommitString)

	baseCommit, bErr := repo.CommitObject(baseHash)
	compareCommit, cErr := repo.CommitObject(compareHash)

	if bErr != nil || cErr != nil {
		logger.Log("Error occurred in fetching commits for the HASH", global.StatusError)
		return res
	}

	baseTree, bTreeErr := baseCommit.Tree()
	compareTree, cTreeErr := compareCommit.Tree()

	if bTreeErr != nil || cTreeErr != nil {
		logger.Log("Error occurred in fetching TREE records for the commits", global.StatusError)
		return res
	}

	logger.Log(fmt.Sprintf("Initiating tree comparison for %s..%s", baseHash.String(), compareHash.String()), global.StatusInfo)
	diff, diffErr := baseTree.Diff(compareTree)
	patchDiff, patchErr := baseTree.Patch(compareTree)

	if diffErr != nil || patchErr != nil {
		logger.Log("Error occurred while comparing the commit Trees", global.StatusError)
		return res
	}

	fileFullName := patchDiff.FilePatches()

	for i, change := range diff {
		action, _ := change.Action()
		file, _ := fileFullName[i].Files()

		actionType := action.String()

		if actionType == "Insert" {
			actionType = "Add"
		}

		if file != nil {
			res = append(res, &model.GitCommitFileResult{
				Type:     actionType[:1],
				FileName: file.Path(),
			})
		} else {
			_, changedFile, _ := change.Files()
			res = append(res, &model.GitCommitFileResult{
				Type:     actionType[:1],
				FileName: changedFile.Name,
			})
		}
	}
	return res
}
