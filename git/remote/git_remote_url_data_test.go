package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type RemoteUrlDataTestSuite struct {
	suite.Suite
	repo          *git2go.Repository
	noHeadRepo    *git2go.Repository
	listRemoteUrl ListRemoteUrl
}

func TestRemoteUrlDataTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteUrlDataTestSuite))
}

func (suite *RemoteUrlDataTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = r
	suite.noHeadRepo = noHeadRepo
	suite.listRemoteUrl = NewRemoteUrlData(r)
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRemotesArePresent_ShouldReturnRemoteUrlList() {
	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.NotZero(len(urlList))
	suite.Equal("https://github.com/neel1996/gitconvex-test.git", *urlList[0])
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRepoIsNil_ShouldReturnNil() {
	suite.listRemoteUrl = NewRemoteUrlData(nil)

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.Nil(urlList)
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRepoHasNoRemotes_ShouldReturnNil() {
	suite.listRemoteUrl = NewRemoteUrlData(suite.noHeadRepo)

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.Nil(urlList)
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRemotesAreNil_ShouldReturnNil() {
	suite.listRemoteUrl = NewRemoteUrlData(&git2go.Repository{
		Remotes: git2go.RemoteCollection{},
	})

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.Nil(urlList)
}
