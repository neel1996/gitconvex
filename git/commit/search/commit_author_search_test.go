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

type CommitAuthorSearchTestSuite struct {
	suite.Suite
	search     Search
	allCommits []git2go.Commit
}

func TestCommitAuthorSearchTestSuite(t *testing.T) {
	suite.Run(t, new(CommitAuthorSearchTestSuite))
}

func (suite *CommitAuthorSearchTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	repo := middleware.NewRepository(r)
	suite.allCommits = test_utils.GetAllTestCommitLogs(repo)
	suite.search = GetSearchAction("author", suite.allCommits)
}

func (suite *CommitAuthorSearchTestSuite) TestSearch_WhenCommitHasMatchingAuthor_ShouldReturnCommits() {
	commits := suite.search.Search("test")

	suite.NotZero(len(commits))
}
