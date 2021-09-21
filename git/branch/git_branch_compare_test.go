package branch

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch/test_utils"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/neel1996/gitconvex/validator"
	mocks2 "github.com/neel1996/gitconvex/validator/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type BranchCompareTestSuite struct {
	suite.Suite
	branchCompare       Compare
	branchValidator     validator.ValidatorWithStringFields
	repo                middleware.Repository
	mockController      *gomock.Controller
	mockRepo            *mocks.MockRepository
	mockBranchValidator *mocks2.MockValidatorWithStringFields
	mockBranch          *mocks.MockBranch
	mockReference       *mocks.MockReference
	mockCommit          *mocks.MockCommit
	baseBranch          string
	compareBranch       string
	testFile            string
}

func TestBranchCompareTestSuite(t *testing.T) {
	suite.Run(t, new(BranchCompareTestSuite))
}

func (suite *BranchCompareTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.branchValidator = validator.NewBranchValidator()
	suite.testFile = os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "compare_test.txt"
	suite.compareBranch = "new_compare"

	addErr := NewAddBranch(middleware.NewRepository(r), validator.NewBranchValidator()).AddBranch(suite.compareBranch, false, nil)
	if addErr != nil {
		fmt.Println(addErr)
	}

	test_utils.StageAndCommitTestFile(suite.repo, suite.compareBranch, suite.testFile)
}

func (suite *BranchCompareTestSuite) SetupTest() {
	suite.baseBranch = "master"
	suite.compareBranch = "new_compare"

	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockBranch = mocks.NewMockBranch(suite.mockController)
	suite.mockReference = mocks.NewMockReference(suite.mockController)
	suite.mockCommit = mocks.NewMockCommit(suite.mockController)
	suite.mockBranchValidator = mocks2.NewMockValidatorWithStringFields(suite.mockController)

	suite.branchCompare = NewBranchCompare(suite.mockRepo, suite.mockBranchValidator)
}

func (suite *BranchCompareTestSuite) TearDownSuite() {
	test_utils.CheckoutTestBranch(suite.repo, suite.baseBranch)

	err := os.Remove(suite.testFile)
	if err != nil {
		return
	}
	test_utils.DeleteTestBranch(suite.repo, suite.compareBranch)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenBranchesHaveDifferentCommits_ShouldReturnDifference() {
	suite.branchCompare = NewBranchCompare(suite.repo, suite.branchValidator)
	compareResults, err := suite.branchCompare.CompareBranch(suite.baseBranch, suite.compareBranch)

	suite.Nil(err)
	suite.NotNil(compareResults)
	suite.NotZero(len(compareResults))
	suite.NotZero(len(compareResults[0].Commits))

	suite.NotNil(*compareResults[0].Commits[0])
	suite.NotEmpty(*compareResults[0].Commits[0].CommitMessage)
	suite.NotEmpty(compareResults[0].Date)
	suite.GreaterOrEqual(len(compareResults[0].Commits), 1)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenBranchValidationFails_ShouldReturnError() {
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.baseBranch, "").Return(errors.New("VALIDATION_ERR"))

	compareResults, err := suite.branchCompare.CompareBranch(suite.baseBranch, "")

	suite.NotNil(err)
	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenBaseBranchLookupFails_ShouldReturnError() {
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.baseBranch, suite.compareBranch).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.baseBranch, git2go.BranchLocal).Return(nil, errors.New("LOOKUP_ERROR"))

	compareResults, err := suite.branchCompare.CompareBranch(suite.baseBranch, suite.compareBranch)

	suite.NotNil(err)
	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenDiffBranchLookupFails_ShouldReturnError() {
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.baseBranch, suite.compareBranch).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.baseBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.compareBranch, git2go.BranchLocal).Return(nil, errors.New("LOOKUP_ERROR"))

	compareResults, err := suite.branchCompare.CompareBranch(suite.baseBranch, suite.compareBranch)

	suite.NotNil(err)
	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenBranchCompareReturnsZero_ShouldReturnError() {
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.baseBranch, suite.compareBranch).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.baseBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.compareBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Reference().Return(suite.mockReference)
	suite.mockBranch.EXPECT().Cmp(suite.mockReference).Return(0)

	compareResults, err := suite.branchCompare.CompareBranch(suite.baseBranch, suite.compareBranch)

	suite.NotNil(err)
	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenCommitLookupForBaseBranchFails_ShouldReturnError() {
	oid, _ := git2go.NewOid("591ee574417890b6771d8c314d6f116586414e29")
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.baseBranch, suite.compareBranch).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.baseBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.compareBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Reference().Return(suite.mockReference)
	suite.mockBranch.EXPECT().Cmp(suite.mockReference).Return(1)
	suite.mockBranch.EXPECT().Target().Return(oid)
	suite.mockRepo.EXPECT().LookupCommitV2(oid).Return(nil, errors.New("LOOKUP_ERROR"))

	compareResults, err := suite.branchCompare.CompareBranch(suite.baseBranch, suite.compareBranch)

	suite.NotNil(err)
	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenCommitLookupForCompareBranchFails_ShouldReturnError() {
	oid, _ := git2go.NewOid("591ee574417890b6771d8c314d6f116586414e29")
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.baseBranch, suite.compareBranch).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.baseBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.compareBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Reference().Return(suite.mockReference)
	suite.mockBranch.EXPECT().Cmp(suite.mockReference).Return(1)
	suite.mockBranch.EXPECT().Target().Return(oid).MaxTimes(2)
	suite.mockRepo.EXPECT().LookupCommitV2(oid).Return(suite.mockCommit, nil)
	suite.mockRepo.EXPECT().LookupCommitV2(oid).Return(nil, errors.New("LOOKUP_ERROR"))

	compareResults, err := suite.branchCompare.CompareBranch(suite.baseBranch, suite.compareBranch)

	suite.NotNil(err)
	suite.Len(compareResults, 0)
}

func (suite *BranchCompareTestSuite) TestCompareBranch_WhenParentCountOfCommitsIsZero_ShouldReturnEmptyResults() {
	oid, _ := git2go.NewOid("591ee574417890b6771d8c314d6f116586414e29")
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.baseBranch, suite.compareBranch).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.baseBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.compareBranch, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Reference().Return(suite.mockReference)
	suite.mockBranch.EXPECT().Cmp(suite.mockReference).Return(1)
	suite.mockBranch.EXPECT().Target().Return(oid).MaxTimes(2)
	suite.mockRepo.EXPECT().LookupCommitV2(oid).Return(suite.mockCommit, nil)
	suite.mockRepo.EXPECT().LookupCommitV2(oid).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().ParentCount().Return(uint(0)).Times(2)

	compareResults, err := suite.branchCompare.CompareBranch(suite.baseBranch, suite.compareBranch)

	suite.Nil(err)
	suite.Len(compareResults, 0)
}
