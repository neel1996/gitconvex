package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/commit"
	"github.com/neel1996/gitconvex/global"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type BranchCompareTestSuite struct {
	suite.Suite
	branchCompare Compare
	repo          *git2go.Repository
	baseBranch    string
	compareBranch string
	testFile      string
}

func TestBranchCompareTestSuite(t *testing.T) {
	suite.Run(t, new(BranchCompareTestSuite))
}

func (suite *BranchCompareTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = r
	suite.testFile = os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "compare_test.txt"
	suite.compareBranch = "new_compare"
	addErr := NewAddBranch(r, suite.compareBranch, false, nil).AddBranch()
	if addErr != nil {
		fmt.Println(addErr)
	}
	suite.stageAndCommitTestFile()
}

func (suite *BranchCompareTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = r
	suite.baseBranch = "master"
	suite.compareBranch = "new_compare"
	suite.branchCompare = NewBranchCompare(suite.repo, suite.baseBranch, suite.compareBranch)
}

func (suite *BranchCompareTestSuite) TearDownSuite() {
	checkoutErr := NewBranchCheckout(suite.repo, "master").CheckoutBranch()
	if checkoutErr != nil {
		logger.Log(checkoutErr.Error(), global.StatusError)
		return
	}

	err := os.Remove(suite.testFile)
	if err != nil {
		return
	}
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenBranchesHaveDifferentCommits_ShouldReturnDifference() {
	compareResults := suite.branchCompare.CompareBranch()

	suite.NotNil(compareResults)
	suite.NotZero(len(compareResults))
	suite.NotZero(len(compareResults[0].Commits))

	suite.NotNil(*compareResults[0].Commits[0])
	suite.NotEmpty(*compareResults[0].Commits[0].CommitMessage)
	suite.NotEmpty(compareResults[0].Date)
	suite.Equal(1, len(compareResults[0].Commits))
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenRepoIsNil_ShouldReturnEmptyDifference() {
	suite.branchCompare = NewBranchCompare(nil, suite.baseBranch, suite.compareBranch)
	compareResults := suite.branchCompare.CompareBranch()

	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenEitherOfBranchInputIsEmpty_ShouldReturnEmptyDifference() {
	suite.branchCompare = NewBranchCompare(suite.repo, "", suite.compareBranch)
	compareResults := suite.branchCompare.CompareBranch()

	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenBranchDoesNotExist_ShouldReturnEmptyDifference() {
	suite.branchCompare = NewBranchCompare(suite.repo, suite.baseBranch, "no_exists")
	compareResults := suite.branchCompare.CompareBranch()

	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenBranchDoesNotDiffer_ShouldReturnEmptyDifference() {
	suite.branchCompare = NewBranchCompare(suite.repo, suite.baseBranch, suite.baseBranch)
	compareResults := suite.branchCompare.CompareBranch()

	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) stageAndCommitTestFile() {
	checkoutErr := NewBranchCheckout(suite.repo, suite.compareBranch).CheckoutBranch()
	if checkoutErr != nil {
		logger.Log(checkoutErr.Error(), global.StatusError)
		return
	}
	err := ioutil.WriteFile(suite.testFile, []byte{0}, 0644)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return
	}
	git.StageItemStruct{
		Repo:     suite.repo,
		FileItem: suite.testFile,
	}.StageItem()

	_ = commit.NewCommitChanges(suite.repo, []string{"Branch compare test commit"}).Add()
}
