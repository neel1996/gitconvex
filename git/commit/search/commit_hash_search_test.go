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

type CommitHashSearchTestSuite struct {
	suite.Suite
	search     Search
	allCommits []git2go.Commit
}

func TestCommitHashSearchTestSuite(t *testing.T) {
	suite.Run(t, new(CommitHashSearchTestSuite))
}

func (suite *CommitHashSearchTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	repo := middleware.NewRepository(r)
	suite.allCommits = test_utils.GetAllTestCommitLogs(repo)
	suite.search = GetSearchAction("hash", suite.allCommits)
}

func (suite *CommitHashSearchTestSuite) TestSearch_ShouldReturnCommitsMatchingHash() {
	commits := suite.search.Search("0")

	suite.NotZero(len(commits))
}
