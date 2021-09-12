package middleware

import git "github.com/libgit2/git2go/v31"

type Commit interface {
	ParentCount() uint
	Id() *git.Oid
	Parent(uint) Commit
	Tree() (*git.Tree, error)
	GetGitCommit() *git.Commit
	Author() *git.Signature
	Message() string
}

type commit struct {
	gitCommit *git.Commit
}

func (c commit) Author() *git.Signature {
	return c.gitCommit.Author()
}

func (c commit) Message() string {
	return c.gitCommit.Message()
}

func (c commit) GetGitCommit() *git.Commit {
	return c.gitCommit
}

func (c commit) Tree() (*git.Tree, error) {
	return c.gitCommit.Tree()
}

func (c commit) Parent(i uint) Commit {
	parent := c.gitCommit.Parent(i)
	if parent == nil {
		return nil
	}

	return NewCommit(parent)
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
