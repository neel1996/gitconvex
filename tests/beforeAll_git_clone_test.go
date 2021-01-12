package tests

import (
	"fmt"
	"github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestCloneHandler(t *testing.T) {
	cwd, _ := os.Getwd()
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
				RepoPath:   tt.args.repoPath,
				RepoURL:    tt.args.repoURL,
				AuthOption: "noauth",
				UserName:   nil,
				Password:   nil,
			}
			got, err := testObject.CloneHandler()
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
