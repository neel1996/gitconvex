package commit

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type CommitChangesTestSuite struct {
	suite.Suite
	repo          *git2go.Repository
	noHeadRepo    *git2go.Repository
	commitChanges Changes
	commitMessage []string
}

func TestCommitChangesTestSuite(t *testing.T) {
	suite.Run(t, new(CommitChangesTestSuite))
}

func (suite *CommitChangesTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head_for_commit"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = r
	suite.noHeadRepo = noHeadRepo

	config, _ := suite.noHeadRepo.Config()
	_ = config.SetString("user.name", "test")
	_ = config.SetString("user.email", "test@test.com")

	suite.commitMessage = []string{"Test commit message " + time.Now().String()}
	suite.commitChanges = NewCommitChanges(suite.repo, suite.commitMessage)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenChangesAreCommitted_ShouldReturnNil() {
	err := suite.commitChanges.Add()

	suite.Nil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenCommitMessageIsMultiLine_ShouldCommitToRepo() {
	suite.commitChanges = NewCommitChanges(suite.repo, []string{"Multi line commit", "Test 1"})

	err := suite.commitChanges.Add()

	suite.Nil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenRepoHasNoSignature_ShouldReturnError() {
	suite.commitChanges = NewCommitChanges(suite.noHeadRepo, []string{})

	config, _ := suite.noHeadRepo.Config()
	_ = config.SetString("user.name", "")
	_ = config.SetString("user.email", "")

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenRepoHasNoHead_ShouldCommitToHead() {
	suite.commitChanges = NewCommitChanges(suite.noHeadRepo, []string{"Initial commit"})

	err := suite.commitChanges.Add()

	suite.Nil(err)
}
