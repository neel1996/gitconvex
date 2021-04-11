package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestAddBranch(t *testing.T) {
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
		want string
	}{
		{name: "Git branch add test case", args: struct {
			repo       *git.Repository
			branchName string
		}{repo: r, branchName: "test"}, want: "BRANCH_CREATION_SUCCESS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj git2.AddBranchInterface
			obj = git2.AddBranchInput{
				Repo:       tt.args.repo,
				BranchName: tt.args.branchName,
			}
			got := obj.AddBranch()
			assert.Equal(t, tt.want, got)

			git2.DeleteBranchInputs{
				Repo:       r,
				BranchName: "test",
			}.DeleteBranch()
		})
	}
}
