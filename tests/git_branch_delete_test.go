package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteBranch(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)

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

			git2.NewAddBranch(r, "test", false, nil).AddBranch()

			got := testObj.DeleteBranch()
			assert.Equal(t, tt.want, got)
		})
	}
}
