package utils

import (
	"flag"
	"github.com/neel1996/gitconvex-server/global"
	"go/types"
	"os"
	"path/filepath"
	"runtime"
)

func DefaultDirSetup() (string, error) {
	var baseDirPath *string
	baseDirFlag := flag.Lookup("basedir")

	if baseDirFlag == nil {
		baseDirPath = flag.String("basedir", "/usr/local/gitconvex", "Default gitconvex directory path for Linux and MacOS")
		flag.Parse()
	} else {
		flagPath := baseDirFlag.Value.String()
		baseDirPath = &flagPath
	}

	if baseDirPath == nil {
		return "", types.Error{Msg: "Unable to parse basedir flag"}
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

	_, statErr := os.Stat(*baseDirPath)
	if statErr != nil {
		err := os.Mkdir(*baseDirPath, 0755)
		if err != nil {
			logger.Log(err.Error(), global.StatusError)
			return "", err
		}
	} else {
		return *baseDirPath, nil
	}

	return "", nil
}
