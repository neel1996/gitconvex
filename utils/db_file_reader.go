package utils

import (
	"encoding/json"
	"github.com/neel1996/gitconvex/global"
	"io/ioutil"
)

type RepoData struct {
	RepoId     string `json:"id"`
	RepoName   string `json:"repoName"`
	RepoPath   string `json:"repoPath"`
	TimeStamp  string `json:"timestamp"`
	AuthOption string `json:"authOption"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
	SSHKeyPath string `json:"sshKeyPath"`
}

// DataStoreFileReader reads the database json file tracking the stored repos and returns the data as a struct
func DataStoreFileReader() []RepoData {
	logger := global.Logger{}

	envConfig := EnvConfigFileReader()
	dbFile := envConfig.DataBaseFile
	dbFileContent, err := ioutil.ReadFile(dbFile)

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
	} else {
		var repoData []RepoData
		if unmarshallErr := json.Unmarshal(dbFileContent, &repoData); unmarshallErr != nil {
			logger.Log(unmarshallErr.Error(), global.StatusError)
		}
		return repoData
	}
	return nil
}
