package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/neel1996/gitconvex-server/utils"
	"sort"
	"strings"
)

// CompareBranch compares two branches and returns the commits which are different from each other
// The function uses the git client to fetch the results as go-git lacks this feature
func CompareBranch(repo *git.Repository, baseBranch string, compareBranch string) []*model.BranchCompareResults {
	var diffCommits []*model.BranchCompareResults
	var commits []model.GitCommits
	w, _ := repo.Worktree()

	if w == nil {
		logger.Log("Unable to fetch repo path", global.StatusError)
		return []*model.BranchCompareResults{}
	} else {
		repoPath := w.Filesystem.Root()
		branchStr := baseBranch + ".." + compareBranch
		args := []string{"log", "--oneline", branchStr, "--pretty=format:%h||%an||%ad||%s", "--date=short"}
		cmd := utils.GetGitClient(repoPath, args)

		cmdStr, cmdErr := cmd.Output()

		if cmdErr != nil {
			logger.Log(fmt.Sprintf("Failed executing command -> %s", cmdErr.Error()), global.StatusError)
			return []*model.BranchCompareResults{}
		}

		outputStr := string(cmdStr)
		splitLines := strings.Split(outputStr, "\n")

		if len(splitLines) == 0 {
			logger.Log("No difference could be arrived between the base and compare branch", global.StatusWarning)
			return []*model.BranchCompareResults{}
		}

		for _, commit := range splitLines {
			splitCommit := strings.Split(commit, "||")
			hash := splitCommit[0]
			author := splitCommit[1]
			timestamp := splitCommit[2]
			message := splitCommit[3]

			commits = append(commits, model.GitCommits{
				Hash:          &hash,
				Author:        &author,
				CommitTime:    &timestamp,
				CommitMessage: &message,
			})
		}

		// Organizing differing commits based on date
		var commitMap = make(map[string][]*model.GitCommits)
		for i := range commits {
			selectedDate := *commits[i].CommitTime

			if commitMap[selectedDate] != nil {
				continue
			}

			j := i
			for j < len(commits) {
				if selectedDate == *commits[j].CommitTime {
					if commitMap[selectedDate] != nil {
						tempCommit := commitMap[selectedDate]
						tempCommit = append(tempCommit, &commits[j])
						commitMap[selectedDate] = tempCommit
					} else {
						tempCommit := []*model.GitCommits{&commits[j]}
						commitMap[selectedDate] = tempCommit
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
}
