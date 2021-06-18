package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"go/types"
	"strings"
)

type List interface {
	ListBranches(branchChan chan ListOfBranches)
}

type listBranch struct {
	repo            *git2go.Repository
	localBranchList []*string
	allBranchList   []*string
}

type ListOfBranches struct {
	CurrentBranch string
	BranchList    []*string
	AllBranchList []*string
}

// ListBranches fetches all the branches from the target repository
// The result will be returned as a struct with the current branch and all the available branches
func (l listBranch) ListBranches(branchChan chan ListOfBranches) {
	var currentBranch string
	repo := l.repo

	defer close(branchChan)

	head, headErr := repo.Head()

	if headErr != nil {
		logger.Log(fmt.Sprintf("Repo head is invalid -> %s", headErr.Error()), global.StatusError)
		nilMsg := []string{"No Branches available", "Repo HEAD is nil"}
		branchChan <- ListOfBranches{
			BranchList:    []*string{&nilMsg[0]},
			CurrentBranch: nilMsg[1],
			AllBranchList: []*string{&nilMsg[0]},
		}
		return
	}

	currentBranch = getCurrentBranchName(head)

	l.allBranchList = append(l.allBranchList, &currentBranch)
	l.localBranchList = append(l.localBranchList, &currentBranch)

	localBranchIterator, itrErr := repo.NewBranchIterator(git2go.BranchAll)

	if itrErr == nil {
		l.runBranchIterator(localBranchIterator, currentBranch)
	}

	branchChan <- ListOfBranches{
		BranchList:    l.localBranchList,
		CurrentBranch: currentBranch,
		AllBranchList: l.allBranchList,
	}
}

func (l listBranch) runBranchIterator(localBranchIterator *git2go.BranchIterator, currentBranch string) {
	_ = localBranchIterator.ForEach(func(branch *git2go.Branch, branchType git2go.BranchType) error {
		branchName, nameErr := branch.Name()
		if nameErr != nil {
			return types.Error{Msg: "branch iterator end"}
		}

		if !branch.IsTag() && !branch.IsNote() && branchName != currentBranch {
			l.classifyRemoteAndLocalBranches(branch, branchName, currentBranch)
		}
		return nil
	})
}

func (l listBranch) classifyRemoteAndLocalBranches(branch *git2go.Branch, branchName string, currentBranch string) {
	if branch.IsRemote() && strings.Contains(branchName, "/") {
		l.getRemoteBranchName(branchName, currentBranch)
	} else {
		l.allBranchList = append(l.allBranchList, &branchName)
		l.localBranchList = append(l.localBranchList, &branchName)
	}
}

func (l listBranch) getRemoteBranchName(branchName string, currentBranch string) {
	splitString := strings.Split(branchName, "/")
	splitBranch := splitString[len(splitString)-1]

	if splitBranch != "HEAD" && splitBranch != currentBranch {
		concatRemote := "remotes/" + strings.Join(splitString, "/")
		l.allBranchList = append(l.allBranchList, &concatRemote)
	}
}

func getCurrentBranchName(head *git2go.Reference) string {
	branch := head.Name()
	splitCurrentBranch := strings.Split(branch, "/")
	branch = splitCurrentBranch[len(splitCurrentBranch)-1]
	return branch
}

func NewBranchListInterface(repo *git2go.Repository) List {
	return listBranch{
		repo: repo,
	}
}
