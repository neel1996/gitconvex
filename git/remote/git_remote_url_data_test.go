package remote

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RemoteUrlDataTestSuite struct {
	suite.Suite
	listRemoteUrl ListRemoteUrl
	validation    Validation
}

func TestRemoteUrlDataTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteUrlDataTestSuite))
}

func (suite *RemoteUrlDataTestSuite) SetupTest() {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	suite.validation = NewRemoteValidation()
	suite.listRemoteUrl = NewRemoteUrlData(r, suite.validation)
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRemotesArePresent_ShouldReturnRemoteUrlList() {
	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.NotZero(len(urlList))
	suite.Equal("https://github.com/neel1996/gitconvex-test.git", *urlList[0])
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRepoIsNil_ShouldReturnNil() {
	suite.listRemoteUrl = NewRemoteUrlData(nil, suite.validation)

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.Nil(urlList)
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRemotesAreNil_ShouldReturnNil() {
	suite.listRemoteUrl = NewRemoteUrlData(
		&git2go.Repository{
			Remotes: git2go.RemoteCollection{},
		},
		suite.validation,
	)

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.Nil(urlList)
}
