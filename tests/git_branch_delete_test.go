package tests

import (
	git2go "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestDeleteBranch(t *testing.T) {
	cwd, _ := os.Getwd()
	r, _ := git2go.OpenRepository(path.Join(cwd, ".."))

	type args struct {
		repo       *git2go.Repository
		branchName string
	}
	tests := []struct {
		name string
		args args
		want *model.BranchDeleteStatus
	}{
		{name: "Git branch deletion test case", args: struct {
			repo       *git2go.Repository
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
			if got := testObj.DeleteBranch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}
