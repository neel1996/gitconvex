package tests

import (
	"fmt"
	"github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestCompareCommit(t *testing.T) {
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

	sampleCommits := git2.CommitLogStruct{
		Repo:            r,
		ReferenceCommit: "",
	}.CommitLogs()

	type args struct {
		repo                *git.Repository
		baseCommitString    string
		compareCommitString string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git commit compare test case", args: struct {
			repo                *git.Repository
			baseCommitString    string
			compareCommitString string
		}{repo: r, baseCommitString: *sampleCommits.Commits[1].Hash, compareCommitString: *sampleCommits.Commits[2].Hash}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.CompareCommitInterface
			testObj = git2.CompareCommitStruct{
				Repo:                tt.args.repo,
				BaseCommitString:    tt.args.baseCommitString,
				CompareCommitString: tt.args.compareCommitString,
			}

			got := testObj.CompareCommit()

			assert.NotZero(t, len(got))
		})
	}
}
