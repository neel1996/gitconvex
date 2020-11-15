package git

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"go/types"
)

// InitHandler initializes a new git repo in the target directory

func InitHandler(repoPath string) (*model.ResponseModel, error) {
	_, err := git.PlainInit(repoPath, false)

	if err != nil {
		logger := &global.Logger{}
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
