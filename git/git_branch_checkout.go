package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
)

// CheckoutBranch checks out the branchName received as argument

func CheckoutBranch(repo *git.Repository, branchName string) string {
	logger := global.Logger{}

	w, _ := repo.Worktree()
	checkoutErr := w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branchName),
		Keep:   true,
	})
	if checkoutErr != nil {
		logger.Log(fmt.Sprintf("Failed to checkout branch - %s --> %v", branchName, checkoutErr.Error()), global.StatusError)
	}
	logger.Log(fmt.Sprintf("Current branch checked out to -> %s", branchName), global.StatusInfo)
	return fmt.Sprintf("Head checked out to branch - %v", branchName)
}
