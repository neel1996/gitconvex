package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type ListRemoteTestSuite struct {
	suite.Suite
	repo       *git2go.Repository
	noHeadRepo *git2go.Repository
	listRemote List
}

func (suite *ListRemoteTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = r
	suite.noHeadRepo = noHeadRepo
	suite.listRemote = NewRemoteList(suite.repo)
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
	suite.listRemote = NewRemoteList(nil)

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Nil(remoteList)
}

func (suite *ListRemoteTestSuite) TestGetAllRemotes_WhenRepoHasNoRemotes_ShouldReturnNil() {
	suite.listRemote = NewRemoteList(suite.noHeadRepo)

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Nil(remoteList)
}
