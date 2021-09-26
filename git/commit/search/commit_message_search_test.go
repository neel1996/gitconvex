package search

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/commit/search/test_utils"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CommitMessageSearchTestSuite struct {
	suite.Suite
	allCommits []git2go.Commit
	search     Search
}

func TestCommitMessageSearchTestSuite(t *testing.T) {
	suite.Run(t, new(CommitMessageSearchTestSuite))
}

func (suite *CommitMessageSearchTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	repo := middleware.NewRepository(r)
	suite.allCommits = test_utils.GetAllTestCommitLogs(repo)
	suite.search = GetSearchAction("message", suite.allCommits)
}

func (suite *CommitMessageSearchTestSuite) TestSearch_ShouldReturnCommitsMatchingTheMessage() {
	commits := suite.search.Search("")

	suite.NotZero(len(commits))
}
