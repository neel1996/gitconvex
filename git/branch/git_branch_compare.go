package branch

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/neel1996/gitconvex/validator"
	"sort"
	"strings"
)

type Compare interface {
	CompareBranch(baseBranchName string, diffBranchName string) ([]*model.BranchCompareResults, error)
}

type branchCompare struct {
	repo            middleware.Repository
	branchValidator validator.ValidatorWithStringFields
}

func logAndReturnError(err error) ([]*model.BranchCompareResults, error) {
	logger.Log(err.Error(), global.StatusError)
	return []*model.BranchCompareResults{}, err
}

// CompareBranch compares two branches and returns the commits which are different from each other
// The function uses the git client to fetch the results as go-git lacks this feature
func (b branchCompare) CompareBranch(baseBranchName string, diffBranchName string) ([]*model.BranchCompareResults, error) {
	var diffCommits []*model.BranchCompareResults
	var filteredCommits []model.GitCommits

	if validationErr := b.branchValidator.ValidateWithFields(baseBranchName, diffBranchName); validationErr != nil {
		return logAndReturnError(validationErr)
	}

	baseBranch, compareBranch, lookupErr := b.lookupBranch(baseBranchName, diffBranchName)
	if lookupErr != nil {
		return logAndReturnError(LookupError)
	}

	compareResult := baseBranch.Cmp(compareBranch.Reference())
	if compareResult == 0 {
		return logAndReturnError(CompareError)
	}

	baseHeadCommit, compareHeadCommit, headCommitErr := b.getHeadCommits(baseBranch, compareBranch)
	if headCommitErr != nil {
		return logAndReturnError(NilHeadError)
	}

	baseBranchCommits, compareBranchCommits := b.getAncestorsOf(baseHeadCommit, compareHeadCommit)

	filteredCommits = filterDifferingCommits(compareBranchCommits, baseBranchCommits)
	commitMap := generateCommitMap(filteredCommits)
	diffCommits = sortCommitsInDescendingOrderOfDate(commitMap, diffCommits)

	return diffCommits, nil
}

func (b branchCompare) lookupBranch(baseBranchName string, diffBranchName string) (middleware.Branch, middleware.Branch, error) {
	var (
		baseBranchErr    error
		compareBranchErr error
	)

	baseBranch, baseBranchErr := b.repo.LookupBranch(baseBranchName, git2go.BranchLocal)

	if baseBranchErr != nil {
		logger.Log(baseBranchErr.Error(), global.StatusError)
		return nil, nil, baseBranchErr
	}

	compareBranch, compareBranchErr := b.repo.LookupBranch(diffBranchName, git2go.BranchLocal)
	if compareBranchErr != nil {
		logger.Log(compareBranchErr.Error(), global.StatusError)
		return nil, nil, compareBranchErr
	}

	return baseBranch, compareBranch, nil
}

func (b branchCompare) getHeadCommits(baseBranch middleware.Branch, compareBranch middleware.Branch) (middleware.Commit, middleware.Commit, error) {
	var (
		baseHeadErr    error
		compareHeadErr error
	)

	baseTarget := baseBranch.Target()
	baseHead, baseHeadErr := b.repo.LookupCommitV2(baseTarget)
	if baseHeadErr != nil {
		logger.Log(baseHeadErr.Error(), global.StatusError)
		return nil, nil, baseHeadErr
	}

	compareTarget := compareBranch.Target()
	compareHead, compareHeadErr := b.repo.LookupCommitV2(compareTarget)
	if compareHeadErr != nil {
		logger.Log(compareHeadErr.Error(), global.StatusError)
		return nil, nil, compareHeadErr
	}

	return baseHead, compareHead, nil
}

func (b branchCompare) getAncestorsOf(baseHead middleware.Commit, compareHead middleware.Commit) (map[string]middleware.Commit, map[string]middleware.Commit) {
	baseBranchMap := getParentCommitsOf(baseHead)
	compareBranchMap := getParentCommitsOf(compareHead)

	return baseBranchMap, compareBranchMap
}

// Organizing differing commits based on date
func generateCommitMap(filteredCommits []model.GitCommits) map[string][]*model.GitCommits {
	var commitMap = make(map[string][]*model.GitCommits)
	var commitHashMap = make(map[string]bool)

	for i := range filteredCommits {
		selectedDate := *filteredCommits[i].CommitTime

		j := i
		for j < len(filteredCommits) {
			if isCommitDateNotTheSelectedDate(selectedDate, filteredCommits, j) && isCommitHashAlreadyInTheList(commitHashMap, filteredCommits, j) {
				break
			}

			if commitMap[selectedDate] != nil {
				tempCommit := commitMap[selectedDate]
				tempCommit = append(tempCommit, &filteredCommits[j])
				commitMap[selectedDate] = tempCommit
				commitHashMap[*filteredCommits[j].Hash] = true
			} else {
				tempCommit := []*model.GitCommits{&filteredCommits[j]}
				commitMap[selectedDate] = tempCommit
				commitHashMap[*filteredCommits[j].Hash] = true
			}
			j++
		}
	}
	return commitMap
}

func isCommitDateNotTheSelectedDate(selectedDate string, filteredCommits []model.GitCommits, iterator int) bool {
	if selectedDate != *filteredCommits[iterator].CommitTime {
		return true
	}

	return false
}

func isCommitHashAlreadyInTheList(commitHashMap map[string]bool, filteredCommits []model.GitCommits, iterator int) bool {
	if len(commitHashMap) > 0 && commitHashMap[*filteredCommits[iterator].Hash] {
		return true
	}

	return false
}

func filterDifferingCommits(compareBranchMap map[string]middleware.Commit, baseBranchMap map[string]middleware.Commit) []model.GitCommits {
	var commits []model.GitCommits
	for commitHash, commit := range compareBranchMap {
		if baseBranchMap[commitHash] == nil {
			hash := commit.Id().String()
			author := commit.Author().Name
			commitTime := strings.Split(commit.Author().When.String(), " ")[0]
			message := commit.Message()

			commits = append(commits, model.GitCommits{
				Hash:          &hash,
				Author:        &author,
				CommitTime:    &commitTime,
				CommitMessage: &message,
			})
		}
	}
	return commits
}

func getParentCommitsOf(head middleware.Commit) map[string]middleware.Commit {
	var next middleware.Commit
	commitMap := make(map[string]middleware.Commit)

	if head.ParentCount() == 0 {
		return commitMap
	}

	next = head.Parent(0)
	for next != nil {
		commitMap[next.Id().String()] = next
		next = next.Parent(0)
	}

	return commitMap
}

func sortCommitsInDescendingOrderOfDate(commitMap map[string][]*model.GitCommits, diffCommits []*model.BranchCompareResults) []*model.BranchCompareResults {
	var keys []string
	for k := range commitMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for _, date := range keys {
		diffCommits = append(diffCommits, &model.BranchCompareResults{
			Date:    date,
			Commits: commitMap[date],
		})
	}
	return diffCommits
}

func NewBranchCompare(repo middleware.Repository, branchValidator validator.ValidatorWithStringFields) Compare {
	return branchCompare{
		repo:            repo,
		branchValidator: branchValidator,
	}
}
