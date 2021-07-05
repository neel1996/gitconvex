package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RemoteAddTestSuite struct {
	suite.Suite
	repo       *git2go.Repository
	remoteName string
	remoteUrl  string
	addRemote  Add
}

func TestRemoteAddTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteAddTestSuite))
}

func (suite *RemoteAddTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = r
	suite.remoteName = "new_origin"
	suite.remoteUrl = "https://github.com/neel1996/gitconvex-test.git"
	suite.addRemote = NewAddRemote(suite.repo, suite.remoteName, suite.remoteUrl)
}

func (suite *RemoteAddTestSuite) TearDownSuite() {
	err := NewDeleteRemote(suite.repo, suite.remoteName).DeleteRemote()
	if err != nil {
		logger.Log(err.Error(), global.StatusWarning)
		return
	}
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenNewRemoteIsAdded_ShouldReturnNoError() {
	err := suite.addRemote.NewRemote()

	suite.Nil(err)
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenRepoIsNil_ShouldReturnError() {
	suite.addRemote = NewAddRemote(nil, suite.remoteName, suite.remoteUrl)

	err := suite.addRemote.NewRemote()

	suite.NotNil(err)
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenRemoteNameIsEmpty_ShouldReturnError() {
	suite.addRemote = NewAddRemote(suite.repo, "", suite.remoteUrl)

	err := suite.addRemote.NewRemote()

	suite.NotNil(err)
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenRemoteUrlIsEmpty_ShouldReturnError() {
	suite.addRemote = NewAddRemote(suite.repo, suite.remoteName, "")

	err := suite.addRemote.NewRemote()

	suite.NotNil(err)
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenRemoteCreationFails_ShouldReturnError() {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))

	suite.addRemote = NewAddRemote(r, "new_origin", "https://github.com/neel1996/gitconvex-test.git")

	err := suite.addRemote.NewRemote()

	suite.NotNil(err)
}
