package api

import (
	"encoding/json"
	"fmt"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/neel1996/gitconvex/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

var logger global.Logger

// GetSettingsData returns the data extracted from the env_config.json file
func GetSettingsData() *model.SettingsDataResults {
	envData := utils.EnvConfigFileReader()

	if envData != nil {
		return &model.SettingsDataResults{
			SettingsDatabasePath: envData.DataBaseFile,
			SettingsPortDetails:  envData.Port,
		}
	} else {
		return &model.SettingsDataResults{}
	}
}

// UpdatePortNumber updates the port number in the env_config.json file with the
// newly supplied port number
func UpdatePortNumber(newPort string) string {
	if utils.EnvConfigValidator() == nil {
		var newEnvData utils.EnvConfig
		envData := utils.EnvConfigFileReader()

		newEnvData.Port = newPort
		newEnvData.DataBaseFile = envData.DataBaseFile

		execName, execErr := os.Executable()
		if execErr != nil {
			logger.Log(execErr.Error(), global.StatusError)
			return global.PortUpdateError
		}

		cwd := filepath.Dir(execName)

		fileString := cwd + "/" + global.EnvFileName
		envContent, _ := json.MarshalIndent(&newEnvData, "", " ")
		writeErr := ioutil.WriteFile(fileString, envContent, 0755)

		if writeErr != nil {
			logger.Log(writeErr.Error(), global.StatusError)
			return global.PortUpdateError
		} else {
			return global.PortUpdateSuccess
		}
	} else {
		return global.PortUpdateError
	}
}

// UpdateDBFilePath function updates the data file in the env_config json file
func UpdateDBFilePath(newFilePath string) string {
	if utils.EnvConfigValidator() == nil {
		var newEnvData utils.EnvConfig
		envData := utils.EnvConfigFileReader()

		newEnvData.Port = envData.Port
		newEnvData.DataBaseFile = newFilePath

		execName, execErr := os.Executable()
		if execErr != nil {
			logger.Log(execErr.Error(), global.StatusError)
			return global.DataFileUpdateError
		}

		cwd := filepath.Dir(execName)

		fileString := cwd + "/" + global.EnvFileName
		envContent, _ := json.MarshalIndent(&newEnvData, "", " ")
		writeErr := ioutil.WriteFile(fileString, envContent, 0755)

		if writeErr != nil {
			logger.Log(writeErr.Error(), global.StatusError)
			return global.DataFileUpdateError
		} else {
			return global.DataFileUpdateSuccess
		}
	} else {
		return global.DataFileUpdateError
	}
}

func reportError(repoId string, errMsg string, errString string) *model.DeleteStatus {
	logger.Log(fmt.Sprintf("%s -> %s", errMsg, errString), global.StatusError)
	return &model.DeleteStatus{
		Status: global.RepoDeleteError,
		RepoID: repoId,
	}
}

// DeleteRepo deletes a repo entry from the env_config json file
// This will just remove the record for gitconvex alone and disturbs the repo in no other way
func DeleteRepo(repoId string) *model.DeleteStatus {
	var repoData []RepoData
	var newRepoData []RepoData
	if utils.EnvConfigValidator() == nil {
		dbFile := utils.EnvConfigFileReader().DataBaseFile

		file, openErr := os.Open(dbFile)

		if openErr != nil {
			return reportError(repoId, "Unable to open DB file ->", openErr.Error())
		}

		fileData, fileErr := ioutil.ReadAll(file)
		if fileErr != nil {
			return reportError(repoId, "Unable to read DB file ->", fileErr.Error())
		}
		unmarshalErr := json.Unmarshal(fileData, &repoData)

		if unmarshalErr != nil {
			return reportError(repoId, "JSON format is incompatible ->", unmarshalErr.Error())
		}

		for _, data := range repoData {
			if data.Id == repoId {
				continue
			}
			newRepoData = append(newRepoData, data)
		}
		dbFileContent, _ := json.MarshalIndent(newRepoData, "", " ")
		if dbFileContent != nil {
			err := ioutil.WriteFile(dbFile, dbFileContent, 0755)
			if err != nil {
				return reportError(repoId, "Failed to update DB file", err.Error())
			}
			return &model.DeleteStatus{
				Status: global.RepoDeleteSuccess,
				RepoID: repoId,
			}
		} else {
			return reportError(repoId, "Failed to update DB file", "")
		}
	} else {
		return reportError(repoId, "Env config file cannot be accessed", "")
	}
}
