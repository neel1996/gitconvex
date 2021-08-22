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
	"strings"
	"testing"
)

type CommitLogMapperTestSuite struct {
	suite.Suite
	mockController  *gomock.Controller
	repo            middleware.Repository
	mockRepo        *mocks.MockRepository
	fileHistory     FileHistory
	mockFileHistory *commitMocks.MockFileHistory
	commits         []git2go.Commit
	mapper          Mapper
}

func TestCommitLogMapperTestSuite(t *testing.T) {
	suite.Run(t, new(CommitLogMapperTestSuite))
}

func (suite *CommitLogMapperTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = middleware.NewRepository(r)

	limit := uint(3)
	commits, _ := NewListAllLogs(suite.repo, &limit, nil).Get()
	suite.commits = commits

	suite.fileHistory = NewFileHistory(suite.repo)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockFileHistory = commitMocks.NewMockFileHistory(suite.mockController)

	suite.mapper = NewMapper(suite.mockFileHistory)
}

func (suite *CommitLogMapperTestSuite) TestMap_ShouldMapCommitLogFieldsToModel() {
	suite.mapper = NewMapper(suite.fileHistory)

	got := suite.mapper.Map(suite.commits)

	wantFileCount, _ := suite.fileHistory.Get(middleware.NewCommit(&suite.commits[0]))

	suite.Len(got, 3)
	suite.Equal(suite.commits[0].Id().String(), *got[0].Hash)
	suite.Equal(suite.commits[0].Author().Name, *got[0].Author)
	suite.Equal(suite.commits[0].Author().When.String(), *got[0].CommitTime)
	suite.Equal(strings.Split(suite.commits[0].Message(), "\n")[0], *got[0].CommitMessage)
	suite.Equal(len(wantFileCount), *got[0].CommitFilesCount)
}

func (suite *CommitLogMapperTestSuite) TestMap_WhenCommitListIsEmpty_ShouldReturnEmptyModelSlice() {
	suite.mapper = NewMapper(suite.fileHistory)

	got := suite.mapper.Map([]git2go.Commit{})

	suite.Len(got, 0)
}

func (suite *CommitLogMapperTestSuite) TestMap_WhenCommitFileHistoryReturnsError_ShouldConsiderFileHistoryAsZero() {
	suite.mapper = NewMapper(suite.mockFileHistory)

	suite.mockFileHistory.EXPECT().Get(middleware.NewCommit(&suite.commits[0])).Return(nil, errors.New("FILE_HISTORY_ERROR")).Times(3)

	got := suite.mapper.Map(suite.commits)

	suite.Len(got, 3)
	suite.Equal(suite.commits[0].Id().String(), *got[0].Hash)
	suite.Equal(suite.commits[0].Author().Name, *got[0].Author)
	suite.Equal(suite.commits[0].Author().When.String(), *got[0].CommitTime)
	suite.Equal(strings.Split(suite.commits[0].Message(), "\n")[0], *got[0].CommitMessage)
	suite.Equal(0, *got[0].CommitFilesCount)
}
