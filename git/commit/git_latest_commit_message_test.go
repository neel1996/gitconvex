package commit

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type LatestMessageTestSuite struct {
	suite.Suite
	mockController *gomock.Controller
	latestMessage  LatestMessage
	repo           middleware.Repository
	noHeadRepo     *git2go.Repository
	mockRepo       *mocks.MockRepository
	mockWalker     *mocks.MockRevWalk
	mockReference  *mocks.MockReference
}

func TestLatestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(LatestMessageTestSuite))
}

func (suite *LatestMessageTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.mockController = gomock.NewController(suite.T())
	suite.repo = middleware.NewRepository(r)

	_ = NewCommitChanges(suite.repo, []string{"New commit"}).Add()

	suite.noHeadRepo = noHeadRepo
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockWalker = mocks.NewMockRevWalk(suite.mockController)
	suite.mockReference = mocks.NewMockReference(suite.mockController)
	suite.latestMessage = NewLatestMessage(suite.mockRepo)
}

func (suite *LatestMessageTestSuite) TestGet_WhenHeadIsValid_ShouldReturnCommitMessage() {
	_ = NewCommitChanges(suite.repo, []string{"New commit"}).Add()

	suite.latestMessage = NewLatestMessage(suite.repo)

	got := suite.latestMessage.Get()

	suite.NotEmpty(got)
	suite.Equal("New commit", got)
}

func (suite *LatestMessageTestSuite) TestGet_WhenCommitMessageIsMultiLine_ShouldReturnFirstLineOfTheMessage() {
	_ = NewCommitChanges(suite.repo, []string{"New commit", "with two lines"}).Add()
	suite.latestMessage = NewLatestMessage(suite.repo)

	got := suite.latestMessage.Get()

	suite.NotEmpty(got)
	suite.Equal("New commit", got)
}

func (suite *LatestMessageTestSuite) TestGet_WhenCommitMessageIsEmpty_ShouldReturnEmptyMessage() {
	_ = NewCommitChanges(suite.repo, []string{""}).Add()
	suite.latestMessage = NewLatestMessage(suite.repo)

	got := suite.latestMessage.Get()

	suite.Empty(got)
}

func (suite *LatestMessageTestSuite) TestGet_WhenHeadIsInValid_ShouldReturnEmptyMessage() {
	suite.repo = middleware.NewRepository(suite.noHeadRepo)
	suite.latestMessage = NewLatestMessage(suite.repo)

	got := suite.latestMessage.Get()

	suite.Empty(got)
}

func (suite *LatestMessageTestSuite) TestGet_WhenCommitLookupFails_ShouldReturnEmptyMessage() {
	oid, _ := git2go.NewOid("test")

	suite.mockRepo.EXPECT().Head().Return(suite.mockReference, nil)
	suite.mockReference.EXPECT().Target().Return(oid)
	suite.mockRepo.EXPECT().LookupCommit(oid).Return(nil, errors.New("LOOKUP_ERROR"))

	got := suite.latestMessage.Get()

	suite.Empty(got)
}
