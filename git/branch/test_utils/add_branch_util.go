package test_utils

import (
	"fmt"
	"github.com/neel1996/gitconvex/git/middleware"
)

func AddTestBranch(repo middleware.Repository, branchName string) {
	head, _ := repo.Head()
	targetCommit, _ := repo.LookupCommit(head.Target())
	fmt.Println(repo.CreateBranch(branchName, targetCommit, false))
}
