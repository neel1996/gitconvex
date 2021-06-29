package tests

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_branchCompare_CompareBranch(t *testing.T) {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))

	type fields struct {
		repo       *git2go.Repository
		baseBranch string
		diffBranch string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Branch compare test", fields: struct {
			repo       *git2go.Repository
			baseBranch string
			diffBranch string
		}{repo: r, baseBranch: "master", diffBranch: "test_compare"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branch.NewAddBranch(r, tt.fields.diffBranch, false, nil).AddBranch()
			branch.NewBranchCheckout(r, tt.fields.diffBranch).CheckoutBranch()
			setupTestRepoStageAndCommit("Comparison commit for branch compare")

			branch.NewBranchCheckout(r, tt.fields.baseBranch).CheckoutBranch()

			b := branch.NewBranchCompare(r, tt.fields.baseBranch, tt.fields.diffBranch)
			got := b.CompareBranch()

			assert.NotEmpty(t, got)
		})
	}
}

func setupTestRepoStageAndCommit(msg string) bool {
	testRepo := os.Getenv("GITCONVEX_TEST_REPO")
	repository, repoErr := git2go.OpenRepository(testRepo)
	if repoErr != nil {
		return true
	}

	initFile := "init_file.txt"
	_ = ioutil.WriteFile(testRepo+"/"+initFile, []byte{byte(63)}, 0755)

	git.StageAllStruct{Repo: repository}.StageAllItems()
	git.CommitStruct{
		Repo:          repository,
		CommitMessage: msg,
		RepoPath:      testRepo,
	}.CommitChanges()

	return false
}
