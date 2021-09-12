package test_utils

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
)

func DeleteTestBranch(repo middleware.Repository, branchName string) {
	branch, _ := repo.LookupBranch(branchName, git2go.BranchLocal)
	fmt.Println(branch.Delete())
}
