package tests

import (
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestDeleteBranch(t *testing.T) {
	cwd, _ := os.Getwd()
	r, _ := git.PlainOpen(path.Join(cwd, ".."))

	type args struct {
		repo       *git.Repository
		branchName string
		forceFlag  bool
	}
	tests := []struct {
		name string
		args args
		want *model.BranchDeleteStatus
	}{
		{name: "Git branch deletion test case", args: struct {
			repo       *git.Repository
			branchName string
			forceFlag  bool
		}{repo: r, branchName: "test", forceFlag: true}, want: &model.BranchDeleteStatus{Status: "BRANCH_DELETE_SUCCESS"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := git2.DeleteBranch(tt.args.repo, tt.args.branchName, tt.args.forceFlag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}
