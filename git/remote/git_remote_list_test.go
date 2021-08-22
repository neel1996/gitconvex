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

type ListRemoteTestSuite struct {
	suite.Suite
	repo                 middleware.Repository
	mockController       *gomock.Controller
	mockRepo             *mocks.MockRepository
	mockRemoteValidation *remoteMocks.MockValidation
	mockRemotes          *remoteMocks.MockRemotes
	listRemote           List
	remoteValidation     Validation
}

func (suite *ListRemoteTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockRemoteValidation = remoteMocks.NewMockValidation(suite.mockController)
	suite.repo = middleware.NewRepository(r)
	suite.remoteValidation = NewRemoteValidation(suite.repo)
	suite.mockRemotes = remoteMocks.NewMockRemotes(suite.mockController)
	suite.listRemote = NewRemoteList(suite.mockRepo, suite.mockRemoteValidation)
}

func TestListRemoteTestSuite(t *testing.T) {
	suite.Run(t, new(ListRemoteTestSuite))
}

func (suite *ListRemoteTestSuite) TestGetAllRemotes_WhenRepoIsValid_ShouldReturnAllRemotes() {
	suite.listRemote = NewRemoteList(suite.repo, suite.remoteValidation)

	expectedRemotes := []*model.RemoteDetails{{
		RemoteName: "origin",
		RemoteURL:  "https://github.com/neel1996/gitconvex-test.git",
	}}

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Equal(len(expectedRemotes), len(remoteList))
	suite.Equal(expectedRemotes[0].RemoteName, remoteList[0].RemoteName)
	suite.Equal(expectedRemotes[0].RemoteURL, remoteList[0].RemoteURL)
}

func (suite *ListRemoteTestSuite) TestGetAllRemotes_WhenRemoteValidationFails_ShouldReturnNil() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(errors.New("VALIDATION_ERROR"))

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Nil(remoteList)
}

func (suite *ListRemoteTestSuite) TestGetAllRemotes_WhenRepoHasNoRemotes_ShouldReturnNil() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(nil)
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes)
	suite.mockRemotes.EXPECT().List().Return([]string{}, errors.New("LIST_ERROR"))

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Nil(remoteList)
}

func (suite *ListRemoteTestSuite) TestGetAllRemotes_WhenRemoteLookupFails_ShouldReturnNil() {
	suite.mockRemoteValidation.EXPECT().ValidateRemoteFields().Return(nil)
	suite.mockRepo.EXPECT().Remotes().Return(suite.mockRemotes).MaxTimes(2)
	suite.mockRemotes.EXPECT().List().Return([]string{"REMOTE"}, nil)
	suite.mockRemotes.EXPECT().Lookup("REMOTE").Return(&git2go.Remote{}, errors.New("LOOKUP_ERROR"))

	remoteList := suite.listRemote.GetAllRemotes()

	suite.Nil(remoteList)
}
