package middleware

import git "github.com/libgit2/git2go/v31"

type Commit interface {
	ParentCount() uint
	Id() *git.Oid
	Parent(uint) Commit
	Tree() (*git.Tree, error)
}

type commit struct {
	gitCommit *git.Commit
}

func (c commit) Tree() (*git.Tree, error) {
	return c.gitCommit.Tree()
}

func (c commit) Parent(i uint) Commit {
	return NewCommit(c.gitCommit.Parent(i))
}

func (c commit) Id() *git.Oid {
	return c.gitCommit.Id()
}

func (c commit) ParentCount() uint {
	return c.gitCommit.ParentCount()
}

func NewCommit(gitCommit *git.Commit) Commit {
	return commit{gitCommit: gitCommit}
}
