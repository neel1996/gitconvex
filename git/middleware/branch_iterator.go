package middleware

import git "github.com/libgit2/git2go/v31"

type BranchIterator interface {
	ForEach(func(b *git.Branch, branchType git.BranchType) error) error
}

type branchIterator struct {
	iterator *git.BranchIterator
}

func (b branchIterator) ForEach(f func(b *git.Branch, branchType git.BranchType) error) error {
	return b.iterator.ForEach(f)
}

func NewBranchIterator(iterator *git.BranchIterator) BranchIterator {
	return branchIterator{iterator: iterator}
}
