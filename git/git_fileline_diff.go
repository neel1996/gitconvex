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
	//var (
	//	currentFileLines []string
	//	commitLines      []string
	//	diffLines        []string
	//	fileDiff         []*string
	//	codedLines       []string
	//	diffIndicator    string
	//	insertionCount   int
	//	deletionCount    int
	//)

	repo := f.Repo
	//data := f.Data
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
		Pathspec: []string{fileName},
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

		fmt.Println(patch.String())
	}

	return returnFileDiffErr("")
	//for _, line := range data {
	//	currentFileLines = append(currentFileLines, *line)
	//}
	//
	//cItr, _ := repo.Log(&git.LogOptions{
	//	All: true,
	//})
	//
	//_ = cItr.ForEach(func(commit *object.Commit) error {
	//	file, _ := commit.File(fileName)
	//	if file != nil {
	//
	//		lines, _ := file.Lines()
	//		commitLines = lines
	//		return types.Error{Msg: "END"}
	//	}
	//	return nil
	//})
	//
	//src := strings.Join(commitLines, "\n")
	//dst := strings.Join(currentFileLines, "\n")
	//
	//diffs := diff.Do(src, dst)
	//
	//diffIndicator = ""
	//
	//for _, d := range diffs {
	//	var indicatedLines []string
	//	splitString := strings.Split(d.Text, "\n")
	//
	//	for _, line := range splitString {
	//		switch d.Type {
	//		case 0:
	//			diffIndicator = ""
	//			break
	//		case 1:
	//			diffIndicator = "+"
	//			insertionCount++
	//			break
	//		case -1:
	//			diffIndicator = "-"
	//			deletionCount++
	//			break
	//		}
	//		if diffIndicator != "" && line == "" {
	//			if diffIndicator == "+" {
	//				insertionCount--
	//			} else {
	//				deletionCount--
	//			}
	//			continue
	//		}
	//		changeStr := diffIndicator + line
	//		indicatedLines = append(indicatedLines, changeStr)
	//	}
	//	codedLines = append(codedLines, indicatedLines...)
	//}
	//
	//for _, line := range codedLines {
	//	diffLines = append(diffLines, line)
	//}
	//
	//diffStat := fmt.Sprintf("%v insertions (+),%v deletions (-)", insertionCount, deletionCount)
	//
	//for i := range diffLines {
	//	fileDiff = append(fileDiff, &diffLines[i])
	//}
	//
	//if insertionCount == 0 && deletionCount == 0 {
	//	msg := "NO_DIFF"
	//	return &model.FileLineChangeResult{
	//		DiffStat: msg,
	//		FileDiff: []*string{&msg},
	//	}
	//}

	//return &model.FileLineChangeResult{
	//	DiffStat: nil,
	//	FileDiff: nil,
	//}
}
