package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

// CommitFileList returns the list of files modified, added or removed as part of a commit
// compares the previous commit tree with the current commit tree and, returns the change type (M|D|A) and file name

func CommitFileList(repo *git.Repository, commitHash string) []*model.GitCommitFileResult {
	logger := global.Logger{}
	var res []*model.GitCommitFileResult

	logger.Log(fmt.Sprintf("Fetching file details for commit %v", commitHash), global.StatusInfo)

	hash := plumbing.NewHash(commitHash)
	commit, err := repo.CommitObject(hash)

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		res = append(res, &model.GitCommitFileResult{
			Type:     "",
			FileName: "",
		})

		return res
	}

	currentTree, _ := commit.Tree()
	prev, parentErr := commit.Parents().Next()

	if parentErr != nil {
		logger.Log(parentErr.Error(), global.StatusError)
		res = append(res, &model.GitCommitFileResult{
			Type:     "",
			FileName: "",
		})

		return res
	}

	prevTree, _ := prev.Tree()
	diff, diffErr := prevTree.Diff(currentTree)

	pDiff, _ := prevTree.Patch(currentTree)
	fileFullName := pDiff.FilePatches()

	if diffErr != nil {
		logger.Log(diffErr.Error(), global.StatusError)
		res = append(res, &model.GitCommitFileResult{
			Type:     "",
			FileName: "",
		})
		return res
	} else {
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

}
