package api

import (
	"encoding/json"
	"fmt"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/neel1996/gitconvex/utils"
	"time"
)

// FetchRepo reads the gitconvex datafile (repo_datastore.json) and returns the stored information
// in the form of categorized array items
func FetchRepo() *model.FetchRepoParams {
	var (
		repoId    []*string
		repoName  []*string
		repoPath  []*string
		timeStamp []*string
	)

	const defaultDateFormat = "2006-01-02 15:04:05"

	repoData := utils.DataStoreFileReader()

	for _, repo := range repoData {
		repoIdStr := repo.RepoId
		repoNameStr := repo.RepoName
		repoPathStr := repo.RepoPath
		timeStampStr := repo.TimeStamp

		convTimeStamp, _ := time.Parse(defaultDateFormat, timeStampStr[:19])
		timeStampStr = convTimeStamp.String()

		repoId = append(repoId, &repoIdStr)
		repoName = append(repoName, &repoNameStr)
		repoPath = append(repoPath, &repoPathStr)
		timeStamp = append(timeStamp, &timeStampStr)
	}

	logger := global.Logger{}
	jsonContent, err := json.MarshalIndent(repoData, "", " ")

	if err != nil {
		logger.Log(fmt.Sprintf("Data file could not be parsed -> %s", err.Error()), global.StatusError)
		return &model.FetchRepoParams{
			RepoID:    nil,
			RepoName:  nil,
			RepoPath:  nil,
			TimeStamp: nil,
		}
	}

	logger.Log(fmt.Sprintf("Available data file content \n%v", string(jsonContent)), global.StatusInfo)
	return &model.FetchRepoParams{
		RepoID:    repoId,
		RepoName:  repoName,
		RepoPath:  repoPath,
		TimeStamp: timeStamp,
	}
}
