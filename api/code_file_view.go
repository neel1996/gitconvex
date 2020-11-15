package api

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"os"
)

// CodeFileView returns the lines from the target file and the latest commit corresponding to the file
func CodeFileView(repo *git.Repository, repoPath string, fileName string) *model.CodeFileType {
	var codeLines []*string

	targetFile := repoPath + "/" + fileName
	logger := global.Logger{}
	file, err := os.Open(targetFile)

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return &model.CodeFileType{
			FileData: nil,
		}
	} else {
		logger.Log(fmt.Sprintf("Reading lines from file --> %s", fileName), global.StatusInfo)
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			codeLines = append(codeLines, &line)
		}
		_ = file.Close()
	}

	return &model.CodeFileType{
		FileData: codeLines,
	}
}
