package tests

import (
	"github.com/neel1996/gitconvex/api"
	"github.com/neel1996/gitconvex/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateRepoName(t *testing.T) {
	var testObj api.AddRepoInterface
	testObj = api.AddRepoInputs{
		RepoName: "test", RepoPath: "..", CloneSwitch: false, RepoURL: "", InitSwitch: false, AuthOption: "", UserName: "", Password: "",
	}
	_ = utils.EnvConfigFileGenerator()
	testRepoId := testObj.AddRepo().RepoID

	type args struct {
		repoId   string
		repoName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "test script for rename repo API", args: struct {
			repoId   string
			repoName string
		}{repoId: testRepoId, repoName: "NewTest"}, want: "Repo name updated successfully", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := api.UpdateRepoName(tt.args.repoId, tt.args.repoName)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
