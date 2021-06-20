package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommitFileList(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)

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
		}{repo: r, commitHash: *sampleCommits.Commits[0].Hash}},
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
