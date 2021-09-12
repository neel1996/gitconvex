package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"go/types"
	"strings"
)

type List interface {
	ListBranches() (model.ListOfBranches, error)
}

type listBranch struct {
	repo *git2go.Repository
}

var (
	localBranchList []string
	allBranchList   []string
)

// ListBranches fetches all the branches from the target repository
// The result will be returned as a struct with the current branch and all the available branches
func (l listBranch) ListBranches() (model.ListOfBranches, error) {
	var currentBranch string
	repo := l.repo

	validationErr := NewBranchFieldsValidation(middleware.NewRepository(repo)).ValidateBranchFields()
	if validationErr != nil {
		logger.Log(validationErr.Error(), global.StatusError)
		return model.ListOfBranches{}, validationErr
	}

	head, headErr := repo.Head()
	if headErr != nil {
		logger.Log(fmt.Sprintf("Repo head is invalid -> %s", headErr.Error()), global.StatusError)
		return model.ListOfBranches{}, headErr
	}

	currentBranch = getCurrentBranchName(head)

	allBranchList = append(allBranchList, currentBranch)
	localBranchList = append(localBranchList, currentBranch)

	localBranchIterator, itrErr := repo.NewBranchIterator(git2go.BranchAll)
	if itrErr != nil {
		logger.Log(itrErr.Error(), global.StatusError)
		return model.ListOfBranches{}, itrErr
	}

	err := l.runBranchIterator(localBranchIterator, currentBranch)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return model.ListOfBranches{}, err
	}

	return model.ListOfBranches{
		CurrentBranch: currentBranch,
		BranchList:    localBranchList,
		AllBranchList: allBranchList,
	}, nil
}

func (l listBranch) runBranchIterator(localBranchIterator *git2go.BranchIterator, currentBranch string) error {
	err := localBranchIterator.ForEach(func(branch *git2go.Branch, branchType git2go.BranchType) error {
		branchName, nameErr := branch.Name()
		if nameErr != nil {
			return types.Error{Msg: "Unable to fetch branch name"}
		}

		if !branch.IsTag() && !branch.IsNote() && branchName != currentBranch {
			l.classifyRemoteAndLocalBranches(branch, branchName, currentBranch)
		}
		return nil
	})

	return err
}

func (l listBranch) classifyRemoteAndLocalBranches(branch *git2go.Branch, branchName string, currentBranch string) {
	if branch.IsRemote() && strings.Contains(branchName, "/") {
		l.getRemoteBranchName(branchName, currentBranch)
	} else {
		allBranchList = append(allBranchList, branchName)
		localBranchList = append(localBranchList, branchName)
	}
}

func (l listBranch) getRemoteBranchName(branchName string, currentBranch string) {
	splitString := strings.Split(branchName, "/")
	splitBranch := splitString[len(splitString)-1]

	if splitBranch != "HEAD" && splitBranch != currentBranch {
		concatRemote := "remotes/" + strings.Join(splitString, "/")
		allBranchList = append(allBranchList, concatRemote)
	}
}

func getCurrentBranchName(head *git2go.Reference) string {
	branch := head.Name()
	splitCurrentBranch := strings.Split(branch, "/")
	branch = splitCurrentBranch[len(splitCurrentBranch)-1]
	return branch
}

func NewBranchList(repo *git2go.Repository) List {
	return listBranch{
		repo: repo,
	}
}
