package test_utils

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
)

func CheckoutTestBranch(repo middleware.Repository, branchName string) {
	branch, _ := repo.LookupBranch(branchName, git2go.BranchLocal)
	targetCommit, _ := repo.LookupCommitV2(branch.Target())
	tree, _ := targetCommit.Tree()
	_ = repo.CheckoutTree(tree, &git2go.CheckoutOptions{
		Strategy:       git2go.CheckoutSafe,
		DisableFilters: false,
	})
	fmt.Println(repo.SetHead("refs/heads/" + branchName))
}
