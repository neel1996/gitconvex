package utils

import (
	"fmt"
	"github.com/neel1996/gitconvex-server/global"
	"os/exec"
	"runtime"
	"strings"
)

func GetGitClient(targetPath string, args []string) exec.Cmd {
	// Checking if git is installed and accessible from command line
	var git string
	var err error
	logger := global.Logger{}
	git, pathErr := exec.LookPath("git")

	if pathErr != nil {
		platform := runtime.GOOS

		logger.Log(fmt.Sprintf("Git is not accessible --> %v", pathErr.Error()), global.StatusWarning)
		if strings.Contains(platform, "windows") {
			git, err = exec.LookPath("./gitclient.exe")
		} else {
			git, err = exec.LookPath("./gitclient")
		}
		if err != nil {
			logger.Log(fmt.Sprintf("Standalone Git client is not available --> %v", err.Error()), global.StatusError)
		} else {
			logger.Log(fmt.Sprintf("Using the standlone git client from %s", git), global.StatusInfo)
		}
	} else {
		logger.Log(fmt.Sprintf("Git is accessible from the commandline --> %v", git), global.StatusInfo)
	}

	cmdArgs := append([]string{git, "-C", targetPath}, args...)

	cmd := exec.Cmd{
		Path: git,
		Args: cmdArgs,
	}
	logger.Log(fmt.Sprintf("Command --> %s", cmd.String()), global.StatusInfo)

	return cmd
}
