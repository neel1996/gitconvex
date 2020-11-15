package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestUnPushedCommits(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
		r, _ = git.PlainOpen(repoPath)
	}

	untrackedResult := "untracked.txt"

	_ = ioutil.WriteFile(untrackedResult, []byte{byte(63)}, 0755)
	git2.StageItem(r, untrackedResult)
	git2.CommitChanges(r, "Test Commit")

	type args struct {
		repo      *git.Repository
		remoteRef string
	}
	tests := []struct {
		name string
		args args
		want []*string
	}{
		{name: "Git unpushed commit list test case", args: struct {
			repo      *git.Repository
			remoteRef string
		}{repo: r, remoteRef: "origin/master"}, want: []*string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := git2.UnPushedCommits(tt.args.repo, tt.args.remoteRef); !strings.Contains(*got[0], "Test Commit") {
				t.Errorf("UnPushedCommits() = %v, want %v", got, tt.want)
			}
		})
	}
}
