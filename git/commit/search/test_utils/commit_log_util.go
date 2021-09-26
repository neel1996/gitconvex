package test_utils

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
)

func GetAllTestCommitLogs(repo middleware.Repository) []git2go.Commit {
	var commits []git2go.Commit

	logItr, _ := repo.Walk()

	err := logItr.PushHead()
	err = logItr.Iterate(func(commit *git2go.Commit) bool {
		if commit != nil {
			commits = append(commits, *commit)
			return true
		}
		return false
	})

	fmt.Println(err)
	return commits
}
