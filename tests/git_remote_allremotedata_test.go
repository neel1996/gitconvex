package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
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

	var r *git.Repository
	cwd, _ := os.Getwd()

	if os.Getenv("GOTESTENV") == "ci" {
		r, _ = git.OpenRepository(path.Join(cwd, ".."))
	} else {
		r, _ = git.OpenRepository(os.Getenv("REPODIR"))
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
			if !assert2.NotZerof(t, got, "Testing non-zero array for GetAllRemotes result") {
				t.Errorf("GetAllRemotes() => Received empty array as result")
			}
		})
	}
}
