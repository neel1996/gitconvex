package utils

import (
	"encoding/json"
	"flag"
	"github.com/neel1996/gitconvex-server/global"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

type EnvConfig struct {
	DataBaseFile string `json:"databaseFile"`
	Port         string `json:"port"`
}

func localLogger(message string, status string) {
	logger := &global.Logger{}
	logger.Log(message, status)
}

// getEnvFilePath returns the default filepath for Gitconvex to store the data file
func getEnvFilePath() (string, error) {
	var baseDirPath string
	baseDirFlag := flag.Lookup("basedir")

	if baseDirFlag != nil {
		baseDirPath = flag.Lookup("basedir").Value.String()
		flag.Parse()

		if runtime.GOOS != "windows" || baseDirPath != "" {
			localLogger("Using default path for data file access -> "+baseDirPath, global.StatusInfo)
			return baseDirPath, nil
		}
	}

	execName, execErr := os.Executable()
	if execErr != nil {
		localLogger(execErr.Error(), global.StatusError)
		return "", execErr
	}

	execPath := filepath.Dir(execName)
	localLogger("Using current exe path for data file access -> "+execPath, global.StatusInfo)
	return execPath, nil
}

// EnvConfigValidator checks if the env_config json file is present and accessible
// If the file is missing or unable to access, then an error will be thrown
func EnvConfigValidator() error {
	cwd, cwdErr := getEnvFilePath()
	if cwdErr != nil {
		return cwdErr
	}

	fileString := cwd + "/" + global.EnvFileName
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
	cwd, cwdErr := getEnvFilePath()
	if cwdErr != nil {
		return nil
	}

	fileString := cwd + "/" + global.EnvFileName
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
	cwd, cwdErr := getEnvFilePath()
	if cwdErr != nil {
		return cwdErr
	}

	fileString := cwd + "/" + global.EnvFileName

	envContent, _ := json.MarshalIndent(&EnvConfig{
		DataBaseFile: cwd + "/" + global.DatabaseFilePath,
		Port:         global.DefaultPort,
	}, "", " ")

	return ioutil.WriteFile(fileString, envContent, 0755)
}
