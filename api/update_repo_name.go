package api

import (
	"fmt"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
)

// UpdateRepoName updates the repository name in the data store JSON file
// with the new file name
func UpdateRepoName(repoId string, repoName string) (string, error) {
	logger.Log(fmt.Sprintf("Initating Repo Name update for %s - to new name %s", repoId, repoName), global.StatusInfo)
	repoDataArray := utils.DataStoreFileReader()
	for idx, repoData := range repoDataArray {
		if repoData.RepoId == repoId {
			logger.Log("Matching repo found in data file", global.StatusInfo)
			repoData.RepoName = repoName
			repoDataArray[idx] = repoData
			break
		}
	}

	logger.Log("Writing updated repo data to the JSON file", global.StatusInfo)
	err := utils.DataFileWriter(repoDataArray)

	if err != nil {
		logger.Log("Repo data update failed --> "+err.Error(), global.StatusError)
		return "Unable to update repo name", err
	} else {
		logger.Log("Repo data updated successfully", global.StatusInfo)
		return "Repo name updated successfully", nil
	}
}
