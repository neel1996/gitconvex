package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestRemoteDataStruct_GetAllRemotes(t *testing.T) {
	type fields struct {
		Repo      *git.Repository
		RemoteURL string
	}

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

	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Test for all remote data", fields: struct {
			Repo      *git.Repository
			RemoteURL string
		}{Repo: r, RemoteURL: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := git2.RemoteDataStruct{
				Repo:      tt.fields.Repo,
				RemoteURL: tt.fields.RemoteURL,
			}
			got := res.GetAllRemotes()

			assert2.NotZero(t, got)
		})
	}
}
