package remote

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	remoteMocks "github.com/neel1996/gitconvex/git/remote/mocks"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RemoteAddTestSuite struct {
	suite.Suite
	mockController       *gomock.Controller
	repo                 middleware.Repository
	remotes              middleware.Remotes
	mockRepo             *mocks.MockRepository
	mockRemotes          *remoteMocks.MockRemotes
	mockRemoteValidation *remoteMocks.MockValidation
	remoteName           string
	remoteUrl            string
	addRemote            Add
	remoteValidation     Validation
}

func TestRemoteAddTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteAddTestSuite))
}

func (suite *RemoteAddTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = middleware.NewRepository(r)
	suite.remotes = middleware.NewRemotes(r.Remotes)
	suite.remoteName = "new_origin"
	suite.remoteUrl = "https://github.com/neel1996/gitconvex-test.git"
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockRemotes = remoteMocks.NewMockRemotes(suite.mockController)
	suite.mockRemoteValidation = remoteMocks.NewMockValidation(suite.mockController)
	suite.remoteValidation = NewRemoteValidation(suite.repo)
	suite.addRemote = NewAddRemote(suite.mockRepo, suite.remoteName, suite.remoteUrl, suite.mockRemoteValidation)
}

func (suite *RemoteAddTestSuite) TearDownSuite() {
	suite.remoteValidation = NewRemoteValidation(suite.repo)

	err := NewDeleteRemote(suite.repo, suite.remoteName, suite.remoteValidation).DeleteRemote()
	if err != nil {
		logger.Log(err.Error(), global.StatusWarning)
		return
	}
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenNewRemoteIsAdded_ShouldReturnNil() {
	suite.addRemote = NewAddRemote(suite.repo, suite.remoteName, suite.remoteUrl, suite.remoteValidation)

	err := suite.addRemote.NewRemote()

	suite.Nil(err)
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenValidationFails_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName, suite.remoteUrl).Return(errors.New("VALIDATION_ERROR"))

	err := suite.addRemote.NewRemote()

	suite.NotNil(err)
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenRemoteCreationFails_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName, suite.remoteUrl).Return(nil)
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes)
	suite.mockRemotes.EXPECT().Create(suite.remoteName, suite.remoteUrl).Return(&git2go.Remote{}, errors.New("REMOTE_ERROR"))

	err := suite.addRemote.NewRemote()

	suite.NotNil(err)
}
