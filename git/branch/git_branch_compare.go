package branch

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"sort"
	"strings"
)

type Compare interface {
	CompareBranch() []*model.BranchCompareResults
}

type branchCompare struct {
	Repo       *git2go.Repository
	BaseBranch string
	DiffBranch string
}

func returnBranchCompareError(errString string) []*model.BranchCompareResults {
	if errString != "" {
		logger.Log(errString, global.StatusWarning)
		return []*model.BranchCompareResults{}
	}
	return nil
}

// CompareBranch compares two branches and returns the commits which are different from each other
// The function uses the git client to fetch the results as go-git lacks this feature
func (b branchCompare) CompareBranch() []*model.BranchCompareResults {
	var diffCommits []*model.BranchCompareResults
	var filteredCommits []model.GitCommits
	repo := b.Repo

	baseBranch, baseBranchErr := repo.LookupBranch(b.BaseBranch, git2go.BranchLocal)
	compareBranch, compareBranchErr := repo.LookupBranch(b.DiffBranch, git2go.BranchLocal)

	if baseBranchErr != nil || compareBranchErr != nil {
		return returnBranchCompareError("Unable to lookup target branches from the repo")
	}

	compareResult := baseBranch.Cmp(compareBranch.Reference)
	if compareResult == 0 {
		return returnBranchCompareError("There are no difference between both the branches")
	}

	baseTarget := baseBranch.Target()
	compareTarget := compareBranch.Target()

	baseHead, _ := repo.LookupCommit(baseTarget)
	compareHead, _ := repo.LookupCommit(compareTarget)

	if baseHead == nil || compareHead == nil {
		return returnBranchCompareError("Branch head is NIL")
	}

	baseBranchMap := make(map[string]*git2go.Commit)
	compareBranchMap := make(map[string]*git2go.Commit)

	baseBranchMap = getParentCommits(baseHead)
	compareBranchMap = getParentCommits(compareHead)

	filteredCommits = filterDifferingCommits(compareBranchMap, baseBranchMap)
	commitMap := generateCommitMaps(filteredCommits)
	diffCommits = sortCommitsBasedOnDate(commitMap, diffCommits)

	return diffCommits
}

// Organizing differing commits based on date
func generateCommitMaps(filteredCommits []model.GitCommits) map[string][]*model.GitCommits {
	var commitMap = make(map[string][]*model.GitCommits)
	var commitHashMap = make(map[string]bool)

	for i := range filteredCommits {
		selectedDate := *filteredCommits[i].CommitTime

		j := i
		for j < len(filteredCommits) {
			if shouldBreakLoop(selectedDate, filteredCommits, j, commitHashMap) {
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

func shouldBreakLoop(selectedDate string, filteredCommits []model.GitCommits, iterator int, commitHashMap map[string]bool) bool {
	if selectedDate != *filteredCommits[iterator].CommitTime {
		return true
	}

	if len(commitHashMap) > 0 && commitHashMap[*filteredCommits[iterator].Hash] {
		return true
	}

	return false
}

func filterDifferingCommits(compareBranchMap map[string]*git2go.Commit, baseBranchMap map[string]*git2go.Commit) []model.GitCommits {
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

func getParentCommits(head *git2go.Commit) map[string]*git2go.Commit {
	var next *git2go.Commit
	commitMap := make(map[string]*git2go.Commit)

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

func sortCommitsBasedOnDate(commitMap map[string][]*model.GitCommits, diffCommits []*model.BranchCompareResults) []*model.BranchCompareResults {
	// Sorting commits in reverse chronological order of date
	var keys []string
	for k := range commitMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	// Extracting differing commits to the resultant model
	for _, date := range keys {
		diffCommits = append(diffCommits, &model.BranchCompareResults{
			Date:    date,
			Commits: commitMap[date],
		})
	}
	return diffCommits
}

func NewBranchCompare(repo *git2go.Repository, baseBranch string, diffBranch string) Compare {
	return branchCompare{
		Repo:       repo,
		BaseBranch: baseBranch,
		DiffBranch: diffBranch,
	}
}
