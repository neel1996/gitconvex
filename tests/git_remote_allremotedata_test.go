package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoteDataStruct_GetAllRemotes(t *testing.T) {
	type fields struct {
		Repo      *git.Repository
		RemoteURL string
	}

	r, _ := git.OpenRepository(TestRepo)

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
