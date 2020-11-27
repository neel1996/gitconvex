package utils

import (
	"encoding/json"
	"github.com/neel1996/gitconvex-server/global"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
)

type EnvConfig struct {
	DataBaseFile string `json:"databaseFile"`
	Port         string `json:"port"`
}

func localLogger(message string, status string) {
	logger := &global.Logger{}
	logger.Log(message, status)
}

// EnvConfigValidator checks if the env_config json file is present and accessible
// If the file is missing or unable to access, then an error will be thrown
func EnvConfigValidator() error {
	execName, execErr := os.Executable()
	cwd := filepath.Dir(execName)

	if execErr != nil {
		localLogger(execErr.Error(), global.StatusError)
		return execErr
	}

	fileString := cwd + "/env_config.json"
	_, openErr := os.Open(fileString)

	if openErr != nil {
		localLogger(openErr.Error(), global.StatusError)
		return openErr
	} else {
		if envContent := EnvConfigFileReader(); envContent == nil {
			localLogger("Unable to read env file", global.StatusError)
			return types.Error{Msg: "Invalid content in env_config file"}
		}
	}
	return nil
}

// EnvConfigFileReader reads the env_config json file and returns the config data as a struct
func EnvConfigFileReader() *EnvConfig {
	execName, execErr := os.Executable()
	cwd := filepath.Dir(execName)

	if execErr != nil {
		localLogger(execErr.Error(), global.StatusError)
		return nil
	}

	fileString := cwd + "/env_config.json"
	envFile, err := os.Open(fileString)

	var envConfig *EnvConfig

	if err != nil {
		localLogger(err.Error(), global.StatusError)
		return nil
	} else {
		if fileContent, openErr := ioutil.ReadAll(envFile); openErr == nil {
			unMarshallErr := json.Unmarshal(fileContent, &envConfig)
			if unMarshallErr == nil {
				return envConfig
			} else {
				localLogger(unMarshallErr.Error(), global.StatusError)
				return nil
			}
		}
	}
	return nil
}

// EnvConfigFileGenerator will be invoked when the EnvConfigValidator returns an error or if EnvConfigFileReader returns no data
// The function generates a new env_config.json file and populates it with the default config data
func EnvConfigFileGenerator() error {
	execName, execErr := os.Executable()
	cwd := filepath.Dir(execName)

	if execErr != nil {
		localLogger(execErr.Error(), global.StatusError)
		return types.Error{Msg: execErr.Error()}
	}

	fileString := cwd + "/env_config.json"

	envContent, _ := json.MarshalIndent(&EnvConfig{
		DataBaseFile: cwd + "/database/repo_datastore.json",
		Port:         "9001",
	}, "", " ")

	return ioutil.WriteFile(fileString, envContent, 0755)
}
