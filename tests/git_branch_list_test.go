package tests

import (
	"github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetBranchList(t *testing.T) {
	b := make(chan branch.ListOfBranches)
	r, _ := git.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))

	type args struct {
		repo       *git.Repository
		branchChan chan branch.ListOfBranches
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git branch list test case", args: struct {
			repo       *git.Repository
			branchChan chan branch.ListOfBranches
		}{repo: r, branchChan: b}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testObj := branch.NewBranchListInterface(tt.args.repo)

			go testObj.ListBranches(tt.args.branchChan)
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
