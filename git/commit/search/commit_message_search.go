package search

import (
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/constants"
	"regexp"
	"strings"
)

type commitMessageSearch struct {
	commits []git.Commit
}

func (m commitMessageSearch) Search(searchKey string) []git.Commit {
	var (
		matchingCommits []git.Commit
		counter         = 0
	)

	for _, commit := range m.commits {
		if m.isExceedingSearchLimit(counter) {
			break
		}

		if isMatch, _ := regexp.MatchString(m.toLower(searchKey), m.toLower(commit.Message())); isMatch {
			matchingCommits = append(matchingCommits, commit)
			counter++
		}
	}

	return matchingCommits
}

func (m commitMessageSearch) toLower(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (m commitMessageSearch) isExceedingSearchLimit(searchLimitCounter int) bool {
	return searchLimitCounter == constants.SearchLimit
}
