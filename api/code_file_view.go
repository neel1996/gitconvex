package api

import (
	"bufio"
	"fmt"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"os"
)

type CodeViewInputs struct {
	RepoPath string
	FileName string
}

type CodeViewInterface interface {
	CodeFileView() *model.CodeFileType
}

// CodeFileView returns the lines from the target file and the latest commit corresponding to the file
func (c CodeViewInputs) CodeFileView() *model.CodeFileType {
	var codeLines []*string

	targetFile := c.RepoPath + "/" + c.FileName
	logger := global.Logger{}
	file, err := os.Open(targetFile)

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return &model.CodeFileType{
			FileData: nil,
		}
	} else {
		logger.Log(fmt.Sprintf("Reading lines from file --> %s", targetFile), global.StatusInfo)
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
