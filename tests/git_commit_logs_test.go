package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCommitLogs(t *testing.T) {
	r, _ := git.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))

	type args struct {
		repo      *git.Repository
		skipCount int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git commit logs test case", args: struct {
			repo      *git.Repository
			skipCount int
		}{repo: r, skipCount: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert2.New(t)

			var testObj git2.CommitLogInterface
			testObj = git2.CommitLogStruct{
				Repo:            tt.args.repo,
				ReferenceCommit: "",
			}

			cLogs := testObj.CommitLogs()
			gotTotal := *cLogs.TotalCommits

			assert.NotZero(int(gotTotal))
		})
	}
}
