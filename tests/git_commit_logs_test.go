package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCommitLogs(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	expectedTotalCommits := float64(18)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
		r, _ = git.PlainOpen(repoPath)
	}

	type args struct {
		repo      *git.Repository
		skipCount int
	}
	tests := []struct {
		name string
		args args
		want *model.GitCommitLogResults
	}{
		{name: "Git commit logs test case", args: struct {
			repo      *git.Repository
			skipCount int
		}{repo: r, skipCount: 0}, want: &model.GitCommitLogResults{
			TotalCommits: &expectedTotalCommits,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert2.New(t)
			cLogs := git2.CommitLogs(tt.args.repo, tt.args.skipCount, "")
			gotTotal := *cLogs.TotalCommits
			assert.Equal(expectedTotalCommits, gotTotal, "Total commit count are mis-matching")
		})
	}
}
