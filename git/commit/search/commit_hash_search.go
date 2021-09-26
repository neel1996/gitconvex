package search

import (
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/constants"
	"regexp"
)

type commitHashSearch struct {
	commits []git.Commit
}

func (h commitHashSearch) Search(searchKey string) []git.Commit {
	var (
		matchingCommits []git.Commit
		counter         = 0
	)

	for _, commit := range h.commits {
		if h.isExceedingSearchLimit(counter) {
			break
		}

		if isMatch, _ := regexp.MatchString(searchKey, commit.Id().String()); isMatch {
			matchingCommits = append(matchingCommits, commit)
			counter++
		}
	}

	return matchingCommits
}

func (h commitHashSearch) isExceedingSearchLimit(searchLimitCounter int) bool {
	return searchLimitCounter == constants.SearchLimit
}
