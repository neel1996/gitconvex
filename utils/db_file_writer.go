package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//DataFileWriter writes a new or updated JSON data to the JSON data file
func DataFileWriter(repoDataArray []RepoData) error {
	envConfig := EnvConfigFileReader()
	dbFile := envConfig.DataBaseFile
	repoDataJSON, _ := json.Marshal(repoDataArray)
	osRead, _ := os.Open(dbFile)
	_, readErr := ioutil.ReadAll(osRead)

	if readErr == nil {
		return ioutil.WriteFile(dbFile, repoDataJSON, 0755)
	} else {
		return readErr
	}
}
