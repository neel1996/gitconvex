package tests

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCommitLogs(t *testing.T) {
	var repoPath string
	var r *git2go.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	expectedTotalCommits := float64(19)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
		r, _ = git2go.OpenRepository(repoPath)
	}

	type args struct {
		repo      *git2go.Repository
		skipCount int
	}
	tests := []struct {
		name string
		args args
		want *model.GitCommitLogResults
	}{
		{name: "Git commit logs test case", args: struct {
			repo      *git2go.Repository
			skipCount int
		}{repo: r, skipCount: 0}, want: &model.GitCommitLogResults{
			TotalCommits: &expectedTotalCommits,
		}},
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
			assert.Equal(expectedTotalCommits, gotTotal, "Total commit count are mis-matching")
		})
	}
}
