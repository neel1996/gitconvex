package commit

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/constants"
	commitMocks "github.com/neel1996/gitconvex/git/commit/mocks"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

type SearchLogsTestSuite struct {
	suite.Suite
	mockController  *gomock.Controller
	mockRepo        *mocks.MockRepository
	mockListLogs    *commitMocks.MockListAllLogs
	mockFileHistory *commitMocks.MockFileHistory
	listAllLogs     ListAllLogs
	mapper          Mapper
	searchLogs      SearchLogs
}

func TestSearchLogsTestSuite(t *testing.T) {
	suite.Run(t, new(SearchLogsTestSuite))
}

func (suite *SearchLogsTestSuite) SetupSuite() {
	r, err := git.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.listAllLogs = NewListAllLogs(middleware.NewRepository(r), nil, nil)
}

func (suite *SearchLogsTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())

	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockListLogs = commitMocks.NewMockListAllLogs(suite.mockController)
	suite.mockFileHistory = commitMocks.NewMockFileHistory(suite.mockController)

	suite.mapper = NewMapper(suite.mockFileHistory)
	suite.searchLogs = NewSearchLogs(suite.mockRepo, suite.mockListLogs, suite.mapper)
}

func (suite *SearchLogsTestSuite) TestGetMatchingLogs_WhenSearchTypeIsHash_ShouldReturnMatchingCommits() {
	commits := suite.commitLogs()
	wantCommit := commits[0]
	searchKey := wantCommit.Id().String()[0:8]

	want := bindWithModel([]git.Commit{wantCommit})

	suite.mockListLogs.EXPECT().Get().Return(commits, nil)
	suite.mockFileHistory.EXPECT().Get(middleware.NewCommit(&wantCommit)).Return(
		[]*model.GitCommitFileResult{{}},
		nil,
	).MaxTimes(constants.SearchLimit)

	got, err := suite.searchLogs.GetMatchingLogs("hash", searchKey)

	suite.Nil(err)
	suite.Equal(*want[0].Hash, *got[0].Hash)
}

func (suite *SearchLogsTestSuite) TestGetMatchingLogs_WhenSearchTypeIsMessage_ShouldReturnMatchingCommits() {
	commits := suite.commitLogs()
	wantCommit := commits[len(commits)-1]
	searchKey := strings.Split(wantCommit.Message(), " ")[0]

	want := bindWithModel([]git.Commit{wantCommit})

	suite.mockListLogs.EXPECT().Get().Return(commits, nil)
	suite.mockFileHistory.EXPECT().Get(middleware.NewCommit(&wantCommit)).Return(
		[]*model.GitCommitFileResult{{}},
		nil,
	).MaxTimes(constants.SearchLimit)

	got, err := suite.searchLogs.GetMatchingLogs("message", searchKey)

	suite.Nil(err)
	suite.Equal(*want[0].CommitMessage, *got[0].CommitMessage)
}

func (suite *SearchLogsTestSuite) TestGetMatchingLogs_WhenSearchTypeIsAuthor_ShouldReturnMatchingCommits() {
	commits := suite.commitLogs()
	wantCommit := commits[0]
	searchKey := wantCommit.Author().Name[0:3]

	want := bindWithModel([]git.Commit{wantCommit})

	suite.mockListLogs.EXPECT().Get().Return(commits, nil)
	suite.mockFileHistory.EXPECT().Get(middleware.NewCommit(&wantCommit)).Return(
		[]*model.GitCommitFileResult{{}},
		nil,
	).MaxTimes(constants.SearchLimit)

	got, err := suite.searchLogs.GetMatchingLogs("author", searchKey)

	suite.Nil(err)
	suite.Equal(*want[0].Author, *got[0].Author)
}

func (suite *SearchLogsTestSuite) TestGetMatchingLogs_WhenSearchTypeIsInvalid_ShouldReturnError() {
	commits := suite.commitLogs()
	wantCommit := commits[0]
	searchKey := wantCommit.Author().Name[0:3]

	suite.mockListLogs.EXPECT().Get().Return(commits, nil)

	_, err := suite.searchLogs.GetMatchingLogs("invalid", searchKey)

	suite.NotNil(err)
	suite.Equal(InvalidSearchCategoryError, err)
}

func (suite *SearchLogsTestSuite) TestGetMatchingLogs_WhenListLogsReturnsError_ShouldReturnError() {
	commits := suite.commitLogs()
	wantCommit := commits[0]
	searchKey := wantCommit.Author().Name[0:3]

	suite.mockListLogs.EXPECT().Get().Return(nil, errors.New("LIST_LOGS_ERROR"))

	_, err := suite.searchLogs.GetMatchingLogs("hash", searchKey)

	suite.NotNil(err)
}

func (suite *SearchLogsTestSuite) TestGetMatchingLogs_WhenCommitLogsExceedsSearchLimit_ShouldReturnMatchingLogsWithinLimit() {
	commits := suite.commitLogs()
	wantCommit := commits[0]
	searchKey := wantCommit.Id().String()[0:8]

	want := bindWithModel([]git.Commit{wantCommit})

	suite.mockListLogs.EXPECT().Get().Return(suite.commitLogsExceedingLimit(wantCommit), nil)
	suite.mockFileHistory.EXPECT().Get(middleware.NewCommit(&wantCommit)).Return(
		[]*model.GitCommitFileResult{{}},
		nil,
	).MaxTimes(constants.SearchLimit)

	got, err := suite.searchLogs.GetMatchingLogs("hash", searchKey)

	suite.Nil(err)
	suite.Equal(*want[0].Hash, *got[0].Hash)
	suite.Len(got, constants.SearchLimit)
}

func bindWithModel(commits []git.Commit) []*model.GitCommits {
	var mappedCommits []*model.GitCommits

	for _, commit := range commits {
		wantedHash := commit.Id().String()
		wantedAuthor := commit.Author().Name
		wantedTime := commit.Author().When.String()
		wantedMessage := strings.TrimSpace(commit.Message())
		wantedFileCount := 1

		mappedCommits = append(mappedCommits, &model.GitCommits{
			Hash:             &wantedHash,
			Author:           &wantedAuthor,
			CommitTime:       &wantedTime,
			CommitMessage:    &wantedMessage,
			CommitFilesCount: &wantedFileCount,
		})
	}

	return mappedCommits
}

func (suite *SearchLogsTestSuite) commitLogs() []git.Commit {
	commits, err := suite.listAllLogs.Get()
	if err != nil {
		return nil
	}

	return commits
}

func (suite *SearchLogsTestSuite) commitLogsExceedingLimit(commit git.Commit) []git.Commit {
	var commits []git.Commit
	limit := constants.SearchLimit * 2

	for i := 0; i < limit; i++ {
		commits = append(commits, commit)
	}

	return commits
}
