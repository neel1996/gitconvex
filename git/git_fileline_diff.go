package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type FileLineDiffInterface interface {
	FileLineDiff() *model.FileLineChangeResult
}

type FileLineDiffStruct struct {
	Repo     *git2go.Repository
	FileName string
	Data     []*string
}

func returnFileDiffErr(msg string) *model.FileLineChangeResult {
	logger.Log(msg, global.StatusError)
	errMsg := "NO_DIFF"
	return &model.FileLineChangeResult{
		DiffStat: errMsg,
		FileDiff: []*string{&errMsg},
	}
}

// FileLineDiff function compares the current version of the target file with the recently comitted version of the file
// and returns the line wise difference. Similar to git diff <filename>
func (f FileLineDiffStruct) FileLineDiff() *model.FileLineChangeResult {
	repo := f.Repo
	fileName := f.FileName

	head, headErr := repo.Head()
	if headErr != nil {
		return returnFileDiffErr(headErr.Error())
	}

	headCommit, headCommitErr := repo.LookupCommit(head.Target())
	if headCommitErr != nil {
		return returnFileDiffErr("Head commit lookup error : " + headCommitErr.Error())
	}

	currentTree, treeErr := headCommit.Tree()
	if treeErr != nil {
		return returnFileDiffErr("Tree lookup error : " + treeErr.Error())
	}

	diff, diffErr := repo.DiffTreeToWorkdir(currentTree, &git2go.DiffOptions{
		Flags:    git2go.DiffSkipBinaryCheck,
		Pathspec: []string{fileName},
		MaxSize:  10000,
	})
	if diffErr != nil {
		return returnFileDiffErr("Diff error : " + diffErr.Error())
	}

	numDeltas := 0
	numDeltas, _ = diff.NumDeltas()
	if numDeltas == 0 {
		return returnFileDiffErr("No delta to compare")
	}

	for i := 0; i < numDeltas; i++ {
		_, deltaErr := diff.Delta(i)
		if deltaErr != nil {
			logger.Log("Delta Error : "+deltaErr.Error(), global.StatusError)
			break
		}
		patch, patchErr := diff.Patch(i)
		if patchErr != nil {
			return returnFileDiffErr("Patch error : " + patchErr.Error())
		}

		diffStats, diffStatsErr := diff.Stats()
		if diffStatsErr != nil {
			logger.Log(diffStatsErr.Error(), global.StatusError)
			msg := "NO_DIFF"
			return &model.FileLineChangeResult{
				DiffStat: msg,
				FileDiff: []*string{&msg},
			}
		}
		fmt.Println(patch.String())

		diffStat := fmt.Sprintf("%v insertions (+),%v deletions (-)", diffStats.Insertions(), diffStats.Deletions())
		return &model.FileLineChangeResult{
			DiffStat: diffStat,
			FileDiff: []*string{&diffStat},
		}
	}
	return returnFileDiffErr("")
}
