package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"io/ioutil"
	"os"
	"testing"
)

func TestCommitChanges(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex/starfleet"
		r, _ = git.PlainOpen(repoPath)
	}

	sampleFile := "untracked.txt"
	err := ioutil.WriteFile(repoPath+"/"+sampleFile, []byte{byte(63)}, 0755)

	fmt.Println(err)
	fmt.Println(git2.StageAllItems(r))

	type args struct {
		repo          *git.Repository
		commitMessage string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Git commit changes test case", args: struct {
			repo          *git.Repository
			commitMessage string
		}{repo: r, commitMessage: "Test commit"}, want: "COMMIT_DONE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := git2.CommitChanges(tt.args.repo, tt.args.commitMessage); got != tt.want {
				t.Errorf("CommitChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
