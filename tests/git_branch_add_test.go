package tests

import (
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddBranch(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)

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
			obj := branch.NewAddBranch(tt.args.repo, tt.args.branchName, false, nil)
			got := obj.AddBranch()
			assert.Equal(t, tt.want, got)

			_, _ = branch.Operation{Delete: branch.NewDeleteBranch(r, "test")}.GitDeleteBranch()
		})
	}
}
