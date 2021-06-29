package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTotalCommitLogs(t *testing.T) {
	r, _ := git.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))

	logChan := make(chan git2.AllCommitData)

	type args struct {
		repo       *git.Repository
		commitChan chan git2.AllCommitData
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git logs test case", args: struct {
			repo       *git.Repository
			commitChan chan git2.AllCommitData
		}{repo: r, commitChan: logChan}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zeroValue := float64(0)

			var testObj git2.AllCommitInterface
			testObj = git2.AllCommitStruct{Repo: tt.args.repo}
			go testObj.AllCommits(tt.args.commitChan)
			commits := <-logChan
			commitLength := commits.TotalCommits

			assert.Greater(t, commitLength, zeroValue)
		})
	}
}
