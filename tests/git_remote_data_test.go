package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestRemoteData(t *testing.T) {
	remoteChan := make(chan git2.RemoteDataModel)
	var repoPath string
	var r *git.Repository

	cwd, _ := os.Getwd()
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = path.Join(cwd, "..")
		r, _ = git.OpenRepository(repoPath)
	} else {
		repoPath = path.Join(cwd, "../..")
		r, _ = git.OpenRepository(repoPath)
	}

	type args struct {
		repo       *git.Repository
		remoteChan chan git2.RemoteDataModel
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git remote data test case", args: struct {
			repo       *git.Repository
			remoteChan chan git2.RemoteDataModel
		}{repo: r, remoteChan: remoteChan}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert2.New(t)
			var testObj git2.RemoteDataInterface
			testObj = git2.RemoteDataStruct{
				Repo: tt.args.repo,
			}
			go testObj.RemoteData(tt.args.remoteChan)
			remoteData := <-remoteChan

			rHost := remoteData.RemoteHost
			rURL := remoteData.RemoteURL

			assert.Contains(*rURL[0], "github")
			assert.Equal(*rHost, "github")
		})
	}
}
