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

func TestCommitFileList(t *testing.T) {
	var repoPath string
	var r *git.Repository
	cwd, _ := os.Getwd()
	currentEnv := os.Getenv("GOTESTENV")
	mockRepoPath := path.Join(cwd, "../..") + "/starfleet"

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
		repo       *git.Repository
		commitHash string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git commit file list test case", args: struct {
			repo       *git.Repository
			commitHash string
		}{repo: r, commitHash: *sampleCommits.Commits[1].Hash}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObject git2.CommitFileListInterface
			testObject = git2.CommitFileListStruct{
				Repo:       tt.args.repo,
				CommitHash: tt.args.commitHash,
			}

			got := testObject.CommitFileList()

			assert.NotZero(t, len(got))
		})
	}
}
