package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearchCommitLogs(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)

	sampleCommits := git2.CommitLogStruct{
		Repo:            r,
		ReferenceCommit: "",
	}.CommitLogs()
	hash := *sampleCommits.Commits[0].Hash

	type args struct {
		repo       *git.Repository
		searchType string
		searchKey  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Git commit log search test case", args: struct {
			repo       *git.Repository
			searchType string
			searchKey  string
		}{repo: r, searchType: "hash", searchKey: hash}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.SearchCommitInterface
			testObj = git2.SearchCommitStruct{
				Repo:       tt.args.repo,
				SearchType: tt.args.searchType,
				SearchKey:  tt.args.searchKey,
			}

			got := testObj.SearchCommitLogs()

			assert.NotZero(t, len(got))
		})
	}
}
