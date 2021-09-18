package branch

import (
	"errors"
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"strings"
)

type List interface {
	ListBranches() (model.ListOfBranches, error)
}

type listBranch struct {
	repo middleware.Repository
}

type branches struct {
	localBranchList  []string
	remoteBranchList []string
	combinedBranches []string
}

func (l listBranch) ListBranches() (model.ListOfBranches, error) {
	var (
		currentBranch   string
		localBranchList []string
		allBranchList   []string
	)
	repo := l.repo

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

	allBranches, err := l.getAllBranches(localBranchIterator, currentBranch)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return model.ListOfBranches{}, err
	}

	allBranchList = append(allBranchList, allBranches.combinedBranches...)
	localBranchList = append(localBranchList, allBranches.localBranchList...)

	return model.ListOfBranches{
		CurrentBranch: currentBranch,
		BranchList:    localBranchList,
		AllBranchList: allBranchList,
	}, nil
}

func (l listBranch) getAllBranches(localBranchIterator middleware.BranchIterator, currentBranch string) (branches, error) {
	var (
		localBranches  []string
		remoteBranches []string
	)

	err := localBranchIterator.ForEach(func(b *git2go.Branch, branchType git2go.BranchType) error {
		branch := middleware.NewBranch(b)
		branchName, nameErr := branch.Name()
		if nameErr != nil {
			return nameErr
		}

		if branch.IsTag() || branch.IsNote() || branchName == currentBranch {
			return nil
		}

		if branch.IsRemote() {
			prefixedRemoteBranchName, err := l.getPrefixedRemoteBranch(branchName, currentBranch)
			if err == nil {
				remoteBranches = append(remoteBranches, prefixedRemoteBranchName)
			}
		} else {
			localBranches = append(localBranches, branchName)
		}

		return nil
	})

	return branches{
		localBranchList:  localBranches,
		remoteBranchList: remoteBranches,
		combinedBranches: append(localBranches, remoteBranches...),
	}, err
}

func (l listBranch) getPrefixedRemoteBranch(branchName string, currentBranch string) (string, error) {
	shortBranchName := l.getShortBranchName(branchName)

	if shortBranchName != "HEAD" && shortBranchName != currentBranch {
		remoteBranchName := "remotes/" + branchName
		return remoteBranchName, nil
	}
	return "", errors.New("invalid remote branch")
}

func (l listBranch) getShortBranchName(b string) string {
	str := strings.Split(b, "/")
	branchName := str[len(str)-1]

	return branchName
}

func getCurrentBranchName(head middleware.Reference) string {
	branch := head.Name()
	splitCurrentBranch := strings.Split(branch, "/")
	branch = splitCurrentBranch[len(splitCurrentBranch)-1]
	return branch
}

func NewBranchList(repo middleware.Repository) List {
	return listBranch{
		repo: repo,
	}
}
