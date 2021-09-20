package branch

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch/test_utils"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type BranchListTestSuite struct {
	suite.Suite
	repo           middleware.Repository
	mockController *gomock.Controller
	mockRepo       *mocks.MockRepository
	mockReference  *mocks.MockReference
	mockIterator   *mocks.MockBranchIterator
	mockBranch     *mocks.MockBranch
	branchList     List
}

func TestBranchListTestSuite(t *testing.T) {
	suite.Run(t, new(BranchListTestSuite))
}

func (suite *BranchListTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)

	test_utils.CheckoutTestBranchWithType(suite.repo, "test_remote", git2go.BranchRemote)
	test_utils.CheckoutTestBranch(suite.repo, "master")
}

func (suite *BranchListTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockReference = mocks.NewMockReference(suite.mockController)
	suite.mockIterator = mocks.NewMockBranchIterator(suite.mockController)
	suite.mockBranch = mocks.NewMockBranch(suite.mockController)
	suite.branchList = NewBranchList(suite.mockRepo)
}

func (suite *BranchListTestSuite) TestListBranches_WhenRepoHasBranches_ShouldReturnBranchList() {
	suite.branchList = NewBranchList(suite.repo)

	branchList, err := suite.branchList.ListBranches()

	log.Println(branchList)

	suite.Nil(err)
	suite.Equal(2, len(branchList.BranchList))
	suite.Equal(3, len(branchList.AllBranchList))
	suite.Equal("master", branchList.CurrentBranch)
}

func (suite *BranchListTestSuite) TestListBranches_WhenRepoHeadIsInvalid_ShouldReturnError() {
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("HEAD_ERROR"))

	_, err := suite.branchList.ListBranches()

	suite.NotNil(err)
}

func (suite *BranchListTestSuite) TestListBranches_WhenNewBranchIteratorFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().Head().Return(suite.mockReference, nil)
	suite.mockReference.EXPECT().Name().Return("refs/head/master")
	suite.mockRepo.EXPECT().NewBranchIterator(git2go.BranchAll).Return(nil, errors.New("ITERATOR_ERR"))

	_, err := suite.branchList.ListBranches()

	suite.NotNil(err)
}

func (suite *BranchListTestSuite) TestListBranches_WhenBranchIteratorReturnsError_ShouldReturnError() {
	suite.mockRepo.EXPECT().Head().Return(suite.mockReference, nil)
	suite.mockReference.EXPECT().Name().Return("refs/head/master")
	suite.mockRepo.EXPECT().NewBranchIterator(git2go.BranchAll).Return(suite.mockIterator, nil)
	suite.mockIterator.EXPECT().ForEach(gomock.Any()).Return(errors.New("iterator error"))

	_, err := suite.branchList.ListBranches()

	suite.NotNil(err)
}
