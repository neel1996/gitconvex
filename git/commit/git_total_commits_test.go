package commit

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	commitMocks "github.com/neel1996/gitconvex/git/commit/mocks"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type TotalCommitsTestSuite struct {
	suite.Suite
	mockController    *gomock.Controller
	total             Total
	repo              middleware.Repository
	mockRepo          *mocks.MockRepository
	mockWalker        *mocks.MockRevWalk
	noHeadRepo        *git2go.Repository
	mockAllCommitLogs *commitMocks.MockListAllLogs
}

func TestTotalCommitsTestSuite(t *testing.T) {
	suite.Run(t, new(TotalCommitsTestSuite))
}

func (suite *TotalCommitsTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.mockController = gomock.NewController(suite.T())
	suite.repo = middleware.NewRepository(r)
	suite.noHeadRepo = noHeadRepo
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockWalker = mocks.NewMockRevWalk(suite.mockController)
	suite.mockAllCommitLogs = commitMocks.NewMockListAllLogs(suite.mockController)
	suite.total = NewTotalCommits(suite.mockAllCommitLogs)
}

func (suite *TotalCommitsTestSuite) TestGet_WhenLogsAreAvailable_ShouldReturnTotal() {
	suite.mockAllCommitLogs.EXPECT().Get().Return([]git2go.Commit{{}, {}}, nil)

	got := suite.total.Get()

	suite.NotZero(got)
}

func (suite *TotalCommitsTestSuite) TestGet_WhenListLogReturnsError_ShouldReturnZero() {
	suite.mockAllCommitLogs.EXPECT().Get().Return(nil, errors.New("LIST_ERROR"))

	got := suite.total.Get()

	suite.Zero(got)
}

func (suite *TotalCommitsTestSuite) TestGet_WhenRepoHasNoCommits_ShouldReturnZero() {
	suite.mockAllCommitLogs.EXPECT().Get().Return([]git2go.Commit{}, nil)

	got := suite.total.Get()

	suite.Zero(got)
}
