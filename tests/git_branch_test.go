package tests

import (
	"fmt"
	"github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestGetBranchList(t *testing.T) {
	b := make(chan git2.Branch)
	var repoPath string
	var r *git.Repository

	cwd, _ := os.Getwd()
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = path.Join(cwd, "..")
		r, _ = git.OpenRepository(repoPath)
	} else {
		repoPath = path.Join(cwd, "../..")
		r, _ = git.OpenRepository(repoPath)
	}

	type args struct {
		repo       *git.Repository
		branchChan chan git2.Branch
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git branch list test case", args: struct {
			repo       *git.Repository
			branchChan chan git2.Branch
		}{repo: r, branchChan: b}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.BranchListInterface
			testObj = git2.BranchListInputs{
				Repo: tt.args.repo,
			}
			go testObj.GetBranchList(tt.args.branchChan)
			branchList := <-tt.args.branchChan

			cBranch := branchList.CurrentBranch
			aBranch := branchList.AllBranchList
			bList := branchList.BranchList

			assert.NotZero(t, len(cBranch))
			assert.NotZero(t, len(aBranch))
			assert.NotZero(t, len(bList))
		})
	}
}
