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

type RemoteEditTestSuite struct {
	suite.Suite
	repo                 middleware.Repository
	mockController       *gomock.Controller
	mockRepo             *mocks.MockRepository
	mockRemoteValidation *remoteMocks.MockValidation
	mockRemoteList       *remoteMocks.MockList
	mockRemotes          *remoteMocks.MockRemotes
	remoteName           string
	remoteUrl            string
	remoteValidation     Validation
	remoteList           List
	editRemote           Edit
}

func TestRemoteEditTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteEditTestSuite))
}

func (suite *RemoteEditTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockRemotes = remoteMocks.NewMockRemotes(suite.mockController)
	suite.mockRemoteList = remoteMocks.NewMockList(suite.mockController)
	suite.mockRemoteValidation = remoteMocks.NewMockValidation(suite.mockController)

	suite.remoteName = "origin"
	suite.remoteUrl = "https://github.com/neel1996/gitconvex-test.git"
	suite.remoteValidation = NewRemoteValidation(suite.repo)
	suite.remoteList = NewRemoteList(suite.repo, suite.remoteValidation)
	suite.editRemote = NewEditRemote(
		suite.mockRepo,
		suite.remoteName,
		suite.remoteUrl,
		suite.mockRemoteValidation,
		suite.mockRemoteList,
	)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteIsEdited_ShouldReturnNil() {
	suite.editRemote = NewEditRemote(suite.repo, suite.remoteName, suite.remoteUrl, suite.remoteValidation, suite.remoteList)

	wantErr := suite.editRemote.EditRemote()

	suite.Nil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteValidationFails_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName, suite.remoteUrl).Return(errors.New("VALIDATION_ERR"))

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteListReturnsError_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName, suite.remoteUrl).Return(nil)
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes)
	suite.mockRemotes.EXPECT().List().Return(nil, errors.New("LIST_ERROR"))

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteIsNotPresent_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName, suite.remoteUrl).Return(nil)
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes)
	suite.mockRemotes.EXPECT().List().Return([]string{}, nil)
	suite.mockRemoteList.EXPECT().GetAllRemotes().Return(nil)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteIsNotPresentInRepo_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName, suite.remoteUrl).Return(nil)
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes)
	suite.mockRemotes.EXPECT().List().Return([]string{}, nil)
	suite.mockRemoteList.EXPECT().GetAllRemotes().Return([]*model.RemoteDetails{{
		RemoteName: "INVALID",
	}})

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteSetUrlFails_ShouldReturnError() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields(suite.remoteName, suite.remoteUrl).Return(nil)
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes)
	suite.mockRemotes.EXPECT().List().Return([]string{}, nil)
	suite.mockRemoteList.EXPECT().GetAllRemotes().Return([]*model.RemoteDetails{{
		RemoteName: suite.remoteName,
		RemoteURL:  suite.remoteUrl,
	}})
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes)
	suite.mockRemotes.EXPECT().SetUrl(suite.remoteName, suite.remoteUrl).Return(errors.New("SET_URL_ERROR"))

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}
