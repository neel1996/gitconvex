package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type RemoteEditTestSuite struct {
	suite.Suite
	repo       *git2go.Repository
	noHeadRepo *git2go.Repository
	remoteName string
	remoteUrl  string
	validation Validation
	editRemote Edit
}

func TestRemoteEditTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteEditTestSuite))
}

func (suite *RemoteEditTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = r
	suite.noHeadRepo = noHeadRepo
	suite.remoteName = "origin"
	suite.remoteUrl = "https://github.com/neel1996/gitconvex-test.git"
	suite.editRemote = NewEditRemote(r, suite.remoteName, suite.remoteUrl)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteIsEdited_ShouldReturnNil() {
	wantErr := suite.editRemote.EditRemote()

	suite.Nil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRepoIsNil_ShouldReturnError() {
	suite.editRemote = NewEditRemote(nil, suite.remoteName, suite.remoteUrl)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteCollectionIsNil_ShouldReturnError() {
	suite.editRemote = NewEditRemote(&git2go.Repository{
		Remotes: git2go.RemoteCollection{},
	}, suite.remoteName, suite.remoteUrl)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteNameIsEmpty_ShouldReturnError() {
	suite.editRemote = NewEditRemote(suite.repo, "", suite.remoteUrl)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteUrlIsEmpty_ShouldReturnError() {
	suite.editRemote = NewEditRemote(suite.repo, suite.remoteName, "")

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRepoHasNoRemotes_ShouldReturnError() {
	suite.editRemote = NewEditRemote(suite.noHeadRepo, suite.remoteName, suite.remoteUrl)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteIsNotPresent_ShouldReturnError() {
	suite.editRemote = NewEditRemote(suite.repo, "no_exists_remote", suite.remoteUrl)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}
