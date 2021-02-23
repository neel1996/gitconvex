package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"os"
	"path"
	"testing"
)

func TestAddRemote(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
		r, _ = git.PlainOpen(repoPath)
	} else {
		cwd, _ := os.Getwd()
		r, _ = git.PlainOpen(path.Join(cwd, ".."))
	}

	type args struct {
		repo       *git.Repository
		remoteName string
		remoteURL  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Git remote add test case", args: struct {
			repo       *git.Repository
			remoteName string
			remoteURL  string
		}{repo: r, remoteName: "github", remoteURL: "https://github.com/neel1996/starfleet.git"}, want: "REMOTE_ADD_SUCCESS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.AddRemoteInterface
			testObj = git2.AddRemoteStruct{
				Repo:       tt.args.repo,
				RemoteName: tt.args.remoteName,
				RemoteURL:  tt.args.remoteURL,
			}
			if got := testObj.AddRemote(); got != tt.want {
				t.Errorf("AddRemote() = %v, want %v", got, tt.want)
			}
		})
	}
}
