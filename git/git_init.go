package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"go/types"
)

// InitHandler initializes a new git repo in the target directory
func InitHandler(repoPath string) (*model.ResponseModel, error) {
	_, err := git2go.InitRepository(repoPath, false)

	if err != nil {
		logger.Log(fmt.Sprintf("Git repo init failed \n%v", err), global.StatusError)
		return nil, types.Error{
			Msg: fmt.Sprintf("Git repo clone failed \n%v", err),
		}
	}

	return &model.ResponseModel{
		Status:    "success",
		Message:   "Git repo has been initialized",
		HasFailed: false,
	}, nil
}
