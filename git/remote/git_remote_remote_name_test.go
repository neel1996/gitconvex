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

type GetRemoteNameTestSuite struct {
	suite.Suite
	getRemoteName        Name
	repo                 middleware.Repository
	remoteValidation     Validation
	remoteList           List
	mockController       *gomock.Controller
	mockRepo             *mocks.MockRepository
	mockRemoteValidation *remoteMocks.MockValidation
	mockRemoteList       *remoteMocks.MockList
	noHeadRepo           *git2go.Repository
}

func (suite *GetRemoteNameTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.remoteValidation = NewRemoteValidation(suite.repo)
	suite.remoteList = NewRemoteList(suite.repo, suite.remoteValidation)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockRemoteValidation = remoteMocks.NewMockValidation(suite.mockController)
	suite.mockRemoteList = remoteMocks.NewMockList(suite.mockController)
	suite.getRemoteName = NewGetRemoteName(
		suite.mockRepo,
		"https://github.com/neel1996/gitconvex-test.git",
		suite.mockRemoteValidation,
		suite.mockRemoteList,
	)
}

func TestGetRemoteNameTestSuite(t *testing.T) {
	suite.Run(t, new(GetRemoteNameTestSuite))
}

func (suite *GetRemoteNameTestSuite) TestGetRemoteName_ShouldReturnRemoteName_WhenRemoteUrlIsValid() {
	suite.getRemoteName = NewGetRemoteName(
		suite.repo,
		"https://github.com/neel1996/gitconvex-test.git",
		suite.remoteValidation,
		suite.remoteList,
	)

	expectedRemote := "origin"

	actualRemote := suite.getRemoteName.GetRemoteNameWithUrl()

	suite.Equal(expectedRemote, actualRemote)
}

func (suite *GetRemoteNameTestSuite) TestGetRemoteName_ShouldReturnEmptyString_WhenValidationFails() {
	expectedRemote := ""

	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(errors.New("VALIDATION_ERROR"))

	actualRemote := suite.getRemoteName.GetRemoteNameWithUrl()

	suite.Equal(expectedRemote, actualRemote)
}

func (suite *GetRemoteNameTestSuite) TestGetRemoteName_ShouldReturnEmptyString_WhenRepoHasNoRemotes() {
	expectedRemote := ""

	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(nil)
	suite.mockRemoteList.EXPECT().GetAllRemotes().Return(nil)

	actualRemote := suite.getRemoteName.GetRemoteNameWithUrl()

	suite.Equal(expectedRemote, actualRemote)
}

func (suite *GetRemoteNameTestSuite) TestGetRemoteName_ShouldReturnEmptyString_WhenRemoteURLIsNotPresentInRepo() {
	expectedRemote := ""

	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(nil)
	suite.mockRemoteList.EXPECT().GetAllRemotes().Return([]*model.RemoteDetails{})

	actualRemote := suite.getRemoteName.GetRemoteNameWithUrl()

	suite.Equal(expectedRemote, actualRemote)
}
