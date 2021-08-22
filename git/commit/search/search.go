package search

import (
	git2go "github.com/libgit2/git2go/v31"
)

type Search interface {
	Search(string) []git2go.Commit
	New([]git2go.Commit) Search
	ToLower(string) string
}
