package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/diff"
	"github.com/neel1996/gitconvex-server/graph/model"
	"go/types"
	"strings"
)

// FileLineDiff function compares the current version of the target file with the recently comitted version of the file
// and returns the line wise difference. Similar to git diff <filename>
func FileLineDiff(repo *git.Repository, fileName string, data []*string) *model.FileLineChangeResult {
	var (
		currentFileLines []string
		commitLines      []string
		diffLines        []string
		fileDiff         []*string
		codedLines       []string
		diffIndicator    string
		insertionCount   int
		deletionCount    int
	)

	for _, line := range data {
		currentFileLines = append(currentFileLines, *line)
	}

	cItr, _ := repo.Log(&git.LogOptions{
		All: true,
	})

	_ = cItr.ForEach(func(commit *object.Commit) error {
		file, _ := commit.File(fileName)
		if file != nil {

			lines, _ := file.Lines()
			commitLines = lines
			return types.Error{Msg: "END"}
		}
		return nil
	})

	src := strings.Join(commitLines, "\n")
	dst := strings.Join(currentFileLines, "\n")

	diffs := diff.Do(src, dst)

	diffIndicator = ""

	for _, d := range diffs {
		var indicatedLines []string
		splitString := strings.Split(d.Text, "\n")

		for _, line := range splitString {
			switch d.Type {
			case 0:
				diffIndicator = ""
				break
			case 1:
				diffIndicator = "+"
				insertionCount++
				break
			case -1:
				diffIndicator = "-"
				deletionCount++
				break
			}
			if diffIndicator != "" && line == "" {
				if diffIndicator == "+" {
					insertionCount--
				} else {
					deletionCount--
				}
				continue
			}
			changeStr := diffIndicator + line
			indicatedLines = append(indicatedLines, changeStr)
		}
		codedLines = append(codedLines, indicatedLines...)
	}

	for _, line := range codedLines {
		diffLines = append(diffLines, line)
	}

	diffStat := fmt.Sprintf("%v insertions (+),%v deletions (-)", insertionCount, deletionCount)

	for i := range diffLines {
		fileDiff = append(fileDiff, &diffLines[i])
	}

	if insertionCount == 0 && deletionCount == 0 {
		msg := "NO_DIFF"
		return &model.FileLineChangeResult{
			DiffStat: msg,
			FileDiff: []*string{&msg},
		}
	}

	return &model.FileLineChangeResult{
		DiffStat: diffStat,
		FileDiff: fileDiff,
	}
}
