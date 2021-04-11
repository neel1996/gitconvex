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

func TestTotalCommitLogs(t *testing.T) {
	var repoPath string
	var r *git.Repository

	cwd, _ := os.Getwd()
	mockRepoPath := path.Join(cwd, "../..") + "/starfleet"
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = mockRepoPath
		r, _ = git.OpenRepository(repoPath)
	} else {
		repoPath = path.Join(cwd, "../..")
		r, _ = git.OpenRepository(repoPath)
	}
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
