package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"io/ioutil"
	"os"
	"testing"
)

func TestResetAllItems(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex/starfleet"
		r, _ = git.PlainOpen(repoPath)
	}

	untrackedResult := "untracked.txt"
	_ = ioutil.WriteFile(untrackedResult, []byte{byte(63)}, 0755)
	git2.StageItem(r, untrackedResult)

	type args struct {
		repo *git.Repository
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Git rest all test case", args: struct{ repo *git.Repository }{repo: r}, want: "STAGE_ALL_REMOVE_SUCCESS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := git2.ResetAllItems(tt.args.repo); got != tt.want {
				t.Errorf("ResetAllItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
