package tests

import (
	"fmt"
	"github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestCloneRepo_ShouldCloneTheTestRepo(t *testing.T) {
	cwd, _ := os.Getwd()
	currentEnv := os.Getenv("GOTESTENV")
	if currentEnv != "ci" {
		t.Skip("Not supported in non-CI mode")
	}

	testRepoPath := path.Join(cwd, "../..") + "/starfleet"
	fmt.Printf("\n\nRepo Path for Testing : %v\n\n", testRepoPath)
	testURL := "https://github.com/neel1996/starfleet.git"

	type args struct {
		repoPath string
		repoURL  string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.ResponseModel
		wantErr bool
	}{
		{name: "Git clone test case", args: struct {
			repoPath string
			repoURL  string
		}{repoPath: testRepoPath, repoURL: testURL}, want: &model.ResponseModel{
			Status:    "success",
			Message:   "Git clone completed",
			HasFailed: false,
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObject git.CloneInterface
			testObject = git.CloneStruct{
				RepoName:   "starfleet",
				RepoPath:   tt.args.repoPath,
				RepoURL:    tt.args.repoURL,
				AuthOption: "noauth",
				SSHKeyPath: "",
				UserName:   "",
				Password:   "",
			}
			got, err := testObject.CloneRepo()
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneRepo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloneRepo_ShouldFailCloneForInvalidRepoURL(t *testing.T) {
	cwd, _ := os.Getwd()
	currentEnv := os.Getenv("GOTESTENV")
	if currentEnv != "ci" {
		t.Skip("Not supported in non-CI mode")
	}

	testRepoPath := path.Join(cwd, "../..") + "/starfleet"
	fmt.Printf("\n\nRepo Path for Testing : %v\n\n", testRepoPath)
	testURL := "https://github.com/neel1996/wrong_repo.git"

	var testObject git.CloneInterface
	testObject = git.CloneStruct{
		RepoName:   "starfleet",
		RepoPath:   testRepoPath,
		RepoURL:    testURL,
		AuthOption: "noauth",
		SSHKeyPath: "",
		UserName:   "",
		Password:   "",
	}

	got, err := testObject.CloneRepo()

	assert.Nil(t, got)
	assert.NotNil(t, err)
}
