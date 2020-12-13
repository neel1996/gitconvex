package tests

import (
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"os"
	"path"
	"testing"
)

func TestAddBranch(t *testing.T) {
	cwd, _ := os.Getwd()
	r, _ := git.PlainOpen(path.Join(cwd, ".."))

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
			if got := git2.AddBranch(tt.args.repo, tt.args.branchName); got != tt.want {
				t.Errorf("AddBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}
