package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type CompareCommitInterface interface {
	CompareCommit() []*model.GitCommitFileResult
}

type CompareCommitStruct struct {
	Repo                *git2go.Repository
	BaseCommitString    string
	CompareCommitString string
}

// CompareCommit compares the tree difference between two commits and returns the file difference status
func (c CompareCommitStruct) CompareCommit() []*model.GitCommitFileResult {
	var res []*model.GitCommitFileResult

	repo := c.Repo
	baseCommitString := c.BaseCommitString
	compareCommitString := c.CompareCommitString

	baseHash, _ := git2go.NewOid(baseCommitString)
	compareHash, _ := git2go.NewOid(compareCommitString)
	if baseHash == nil || compareHash == nil {
		logger.Log("Converting commit hash string to Oid", global.StatusError)
		return res
	}

	baseCommit, bErr := repo.LookupCommit(baseHash)
	compareCommit, cErr := repo.LookupCommit(compareHash)

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
	diff, diffErr := repo.DiffTreeToTree(baseTree, compareTree, nil)

	if diffErr != nil {
		logger.Log("Error occurred while comparing the commit Trees", global.StatusError)
		logger.Log(diffErr.Error(), global.StatusError)
		return res
	}

	numDeltas := 0
	numDeltas, _ = diff.NumDeltas()

	if numDeltas == 0 {
		logger.Log("No changes found between the trees", global.StatusWarning)
		return res
	}

	for i := 0; i < numDeltas; i++ {
		delta, deltaErr := diff.Delta(i)
		if deltaErr != nil {
			logger.Log(deltaErr.Error(), global.StatusError)
			return res
		}
		filePath := delta.NewFile
		status := delta.Status.String()
		res = append(res, &model.GitCommitFileResult{
			Type:     status[0:1],
			FileName: filePath.Path,
		})
	}

	return res
}
