package api

import (
	"fmt"
	"github.com/neel1996/gitconvex-server/utils"
	"runtime"
	"strings"

	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

func getOs() string {
	return runtime.GOOS
}

func getGitVersion() string {
	gitCmd := utils.GetGitClient(".", []string{"version"})
	gitVersion, err := gitCmd.Output()

	if err != nil {
		fmt.Printf("Git version could not be obtained \n %v", err)
	}

	return strings.Split(string(gitVersion), "\n")[0]
}

// HealthCheckApi returns the current version of git installed in the host and the platform
// gitconvex is running on
func HealthCheckApi() *model.HealthCheckParams {

	logger := global.Logger{Message: fmt.Sprintf("Obtained host information : %v -- %v", getOs(), getGitVersion())}
	logger.LogInfo()

	return &model.HealthCheckParams{
		Os:  getOs(),
		Git: getGitVersion(),
	}
}
