package remote

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	remoteMocks "github.com/neel1996/gitconvex/git/remote/mocks"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RemoteDeleteTestSuite struct {
	suite.Suite
	repo                 middleware.Repository
	remoteName           string
	deleteRemote         Delete
	remoteValidation     Validation
	mockController       *gomock.Controller
	mockRepo             *mocks.MockRepository
	mockRemoteValidation *remoteMocks.MockValidation
	mockRemotes          *remoteMocks.MockRemotes
}

func TestRemoteDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteDeleteTestSuite))
}

func (suite *RemoteDeleteTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = middleware.NewRepository(r)
	suite.remoteName = "new_origin"
	suite.remoteValidation = NewRemoteValidation(suite.repo)

	_ = NewAddRemote(suite.repo, suite.remoteName, "remote://some_url", suite.remoteValidation).NewRemote()
}

func (suite *RemoteDeleteTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = middleware.NewRepository(r)
	suite.remoteName = "new_origin"
	suite.remoteValidation = NewRemoteValidation(suite.repo)

	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockRemoteValidation = remoteMocks.NewMockValidation(suite.mockController)
	suite.mockRemotes = remoteMocks.NewMockRemotes(suite.mockController)
	suite.deleteRemote = NewDeleteRemote(suite.mockRepo, suite.remoteName, suite.mockRemoteValidation)
}

func (suite *RemoteDeleteTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *RemoteDeleteTestSuite) TestDeleteNewRemote_WhenNewRemoteIsDeleted_ShouldReturnNoError() {
	suite.deleteRemote = NewDeleteRemote(suite.repo, suite.remoteName, suite.remoteValidation)
	err := suite.deleteRemote.DeleteRemote()

	suite.Nil(err)
}

func (suite *RemoteDeleteTestSuite) TestDeleteNewRemote_WhenRemoteValidationFails_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName).Return(errors.New("VALIDATION_ERROR"))

	err := suite.deleteRemote.DeleteRemote()

	suite.NotNil(err)
}

func (suite *RemoteDeleteTestSuite) TestDeleteNewRemote_WhenRemoteDeletionFails_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName).Return(nil)
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes)
	suite.mockRemotes.EXPECT().Delete(suite.remoteName).Return(errors.New("DELETION_FAILS"))

	err := suite.deleteRemote.DeleteRemote()

	suite.NotNil(err)
}
