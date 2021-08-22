package remote

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	remoteMocks "github.com/neel1996/gitconvex/git/remote/mocks"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RemoteUrlDataTestSuite struct {
	suite.Suite
	repo                 middleware.Repository
	mockController       *gomock.Controller
	mockRepo             *mocks.MockRepository
	mockRemoteList       *remoteMocks.MockList
	mockRemoteValidation *remoteMocks.MockValidation
	remoteValidation     Validation
	remoteList           List
	listRemoteUrl        ListRemoteUrl
}

func TestRemoteUrlDataTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteUrlDataTestSuite))
}

func (suite *RemoteUrlDataTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.remoteValidation = NewRemoteValidation(suite.repo)
	suite.remoteList = NewRemoteList(suite.repo, suite.remoteValidation)
	suite.remoteValidation = NewRemoteValidation(suite.repo)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockRemoteValidation = remoteMocks.NewMockValidation(suite.mockController)
	suite.mockRemoteList = remoteMocks.NewMockList(suite.mockController)
	suite.listRemoteUrl = NewRemoteUrlData(suite.mockRepo, suite.mockRemoteValidation, suite.mockRemoteList)
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRemotesArePresent_ShouldReturnRemoteUrlList() {
	suite.listRemoteUrl = NewRemoteUrlData(suite.repo, suite.remoteValidation, suite.remoteList)

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.NotZero(len(urlList))
	suite.Equal("https://github.com/neel1996/gitconvex-test.git", *urlList[0])
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRemoteValidationFails_ShouldReturnNil() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(errors.New("VALIDATION_ERR"))

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.Nil(urlList)
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRepoHasNoRemotes_ShouldReturnNil() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(nil)
	suite.mockRemoteList.EXPECT().GetAllRemotes().Return(nil)

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.Nil(urlList)
}

func (suite *RemoteUrlDataTestSuite) TestGetAllRemoteUrl_WhenRemotesAreNil_ShouldReturnNil() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(nil)
	suite.mockRemoteList.EXPECT().GetAllRemotes().Return([]*model.RemoteDetails{})

	urlList := suite.listRemoteUrl.GetAllRemoteUrl()

	suite.Nil(urlList)
}
