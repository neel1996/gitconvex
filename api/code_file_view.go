package api

import (
	"bufio"
	"fmt"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
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
	file, err := os.Open(targetFile)

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return &model.CodeFileType{
			FileData: nil,
		}
	} else {
		fileInfo, fileInfoErr := file.Stat()
		if fileInfoErr != nil {
			logger.Log(fileInfoErr.Error(), global.StatusError)
			return &model.CodeFileType{
				FileData: nil,
			}
		} else {
			if fileInfo.IsDir() {
				logger.Log("Directory cannot be read as a file!", global.StatusError)
				return &model.CodeFileType{
					FileData: nil,
				}
			}
		}

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
