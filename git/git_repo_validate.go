package git

import (
	"fmt"
	"go/types"
	"os"
)

// RepoValidator validates if the specified repoPath represents a valid git repo

func RepoValidator(repoPath string) (string, error) {
	_, err := os.Open(fmt.Sprintf("%v/.git/", repoPath))
	if err != nil {
		return "", types.Error{Msg: "The selected folder is not a git repo"}
	}
	return "Repo is valid!", nil
}
