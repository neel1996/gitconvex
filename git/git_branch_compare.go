package git

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"sort"
	"strings"
)

type BranchCompareInterface interface {
	CompareBranch() []*model.BranchCompareResults
}

type BranchCompareInputs struct {
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
func (inputs BranchCompareInputs) CompareBranch() []*model.BranchCompareResults {
	var diffCommits []*model.BranchCompareResults
	var commits []model.GitCommits
	repo := inputs.Repo

	baseBranch, baseBranchErr := repo.LookupBranch(inputs.BaseBranch, git2go.BranchLocal)
	compareBranch, compareBranchErr := repo.LookupBranch(inputs.DiffBranch, git2go.BranchLocal)

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

	var baseNext *git2go.Commit
	var compareNext *git2go.Commit
	baseBranchMap := make(map[string]*git2go.Commit)
	compareBranchMap := make(map[string]*git2go.Commit)

	if baseHead.ParentCount() > 0 {
		baseNext = baseHead.Parent(0)
		for baseNext != nil {
			baseBranchMap[baseNext.Id().String()] = baseNext
			baseNext = baseNext.Parent(0)
		}
	}

	if compareHead.ParentCount() > 0 {
		compareNext = compareHead.Parent(0)
		for compareNext != nil {
			compareBranchMap[compareNext.Id().String()] = compareNext
			compareNext = compareNext.Parent(0)
		}
	}

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

	// Organizing differing commits based on date
	var commitMap = make(map[string][]*model.GitCommits)
	var tempCommitHashMap = make(map[string]bool)

	for i := range commits {
		selectedDate := *commits[i].CommitTime

		j := i
		for j < len(commits) {
			if selectedDate == *commits[j].CommitTime {
				if len(tempCommitHashMap) > 0 && tempCommitHashMap[*commits[j].Hash] {
					break
				}
				if commitMap[selectedDate] != nil {
					tempCommit := commitMap[selectedDate]
					tempCommit = append(tempCommit, &commits[j])
					commitMap[selectedDate] = tempCommit
					tempCommitHashMap[*commits[j].Hash] = true
				} else {
					tempCommit := []*model.GitCommits{&commits[j]}
					commitMap[selectedDate] = tempCommit
					tempCommitHashMap[*commits[j].Hash] = true
				}
			} else {
				break
			}
			j++
		}
	}

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
