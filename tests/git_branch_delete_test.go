package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestDeleteBranch(t *testing.T) {
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
		branchName string
	}
	tests := []struct {
		name string
		args args
		want *model.BranchDeleteStatus
	}{
		{name: "Git branch deletion test case", args: struct {
			repo       *git.Repository
			branchName string
		}{repo: r, branchName: "test"}, want: &model.BranchDeleteStatus{Status: "BRANCH_DELETE_SUCCESS"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.DeleteBranchInterface
			testObj = git2.DeleteBranchInputs{
				Repo:       tt.args.repo,
				BranchName: tt.args.branchName,
			}

			_ = git2.AddBranchInput{
				Repo:       r,
				BranchName: "test",
			}.AddBranch()

			got := testObj.DeleteBranch()
			assert.Equal(t, tt.want, got)
		})
	}
}
