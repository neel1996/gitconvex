package api

import (
	"encoding/json"
	"fmt"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/neel1996/gitconvex-server/utils"
	"time"
)

func FetchRepo() *model.FetchRepoParams {
	var (
		repoId    []*string
		repoName  []*string
		repoPath  []*string
		timeStamp []*string
	)

	repoData := utils.DataStoreFileReader()

	for _, repo := range repoData {
		repoIdStr := repo.RepoId
		repoNameStr := repo.RepoName
		repoPathStr := repo.RepoPath
		timeStampStr := repo.TimeStamp

		convTimeStamp, _ := time.Parse("2006-01-02 15:04:05", timeStampStr[:19])
		timeStampStr = convTimeStamp.String()

		repoId = append(repoId, &repoIdStr)
		repoName = append(repoName, &repoNameStr)
		repoPath = append(repoPath, &repoPathStr)
		timeStamp = append(timeStamp, &timeStampStr)
	}

	logger := global.Logger{}
	jsonContent, _ := json.MarshalIndent(repoData, "", " ")
	logger.Log(fmt.Sprintf("Reading data file content \n%v", string(jsonContent)), global.StatusInfo)

	return &model.FetchRepoParams{
		RepoID:    repoId,
		RepoName:  repoName,
		RepoPath:  repoPath,
		TimeStamp: timeStamp,
	}
}
