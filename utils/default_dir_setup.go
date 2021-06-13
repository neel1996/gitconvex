package utils

import (
	"github.com/neel1996/gitconvex/global"
	"go/types"
	"os"
	"path/filepath"
	"runtime"
)

func DefaultDirSetup() (string, error) {
	baseDirPath := global.GetDefaultDirectory()

	if baseDirPath == "" {
		return "", types.Error{Msg: "Default directory path value is empty"}
	}

	if runtime.GOOS == "windows" {
		execName, execErr := os.Executable()
		if execErr != nil {
			logger.Log(execErr.Error(), global.StatusError)
			return "", execErr
		}

		execPath := filepath.Dir(execName)
		logger.Log("Using current exe path for data file access -> "+execPath, global.StatusInfo)

		return execPath, nil
	}

	_, statErr := os.Stat(baseDirPath)
	if statErr != nil {
		err := os.Mkdir(baseDirPath, 0644)
		if err != nil {
			logger.Log(err.Error(), global.StatusError)
			return "", err
		}
	} else {
		return baseDirPath, nil
	}

	return "", nil
}
