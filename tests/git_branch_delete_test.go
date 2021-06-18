package tests

import (
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch"
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
			b := branch.Operation{
				Delete: branch.NewDeleteBranch(r, "test"),
				Add:    branch.NewAddBranch(r, "test", false, nil),
			}

			_, err := b.GitAddBranch()
			assert.Nil(t, err)

			got, err := b.GitDeleteBranch()

			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
