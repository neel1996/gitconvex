package remote

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ListRemoteTestSuite struct {
	suite.Suite
	listRemote List
	validate   Validation
}

func (suite *ListRemoteTestSuite) SetupTest() {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	suite.validate = NewRemoteValidation()
	suite.listRemote = NewRemoteList(r, suite.validate)
}

func TestListRemoteTestSuite(t *testing.T) {
	suite.Run(t, new(ListRemoteTestSuite))
}

func (suite *ListRemoteTestSuite) TestGetAllRemotes_WhenRepoIsValid_ShouldReturnAllRemotes() {
	expectedRemotes := []*model.RemoteDetails{{
		RemoteName: "origin",
		RemoteURL:  "https://github.com/neel1996/gitconvex-test.git",
	}}

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Equal(len(expectedRemotes), len(remoteList))
	suite.Equal(expectedRemotes[0].RemoteName, remoteList[0].RemoteName)
	suite.Equal(expectedRemotes[0].RemoteURL, remoteList[0].RemoteURL)
}

func (suite *ListRemoteTestSuite) TestGetAllRemotes_WhenRepoIsNil_ShouldReturnNil() {
	suite.listRemote = NewRemoteList(nil, suite.validate)

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Nil(remoteList)
}

func (suite *ListRemoteTestSuite) TestGetAllRemotes_WhenRepoHasNoRemotes_ShouldReturnNil() {
	r := &git2go.Repository{
		Remotes: git2go.RemoteCollection{},
	}

	suite.listRemote = NewRemoteList(r, suite.validate)

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Nil(remoteList)
}
