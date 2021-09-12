package middleware

import git2go "github.com/libgit2/git2go/v31"

type Branch interface {
	Target() *git2go.Oid
	Delete() error
	Cmp(r Reference) int
	Reference() Reference
}

type branch struct {
	branch *git2go.Branch
}

func (b branch) Reference() Reference {
	return NewReference(b.branch.Reference)
}

func (b branch) Cmp(r Reference) int {
	return b.branch.Cmp(r.GetGitReference())
}

func (b branch) Delete() error {
	return b.branch.Delete()
}

func (b branch) Target() *git2go.Oid {
	return b.branch.Target()
}

func NewBranch(gitBranch *git2go.Branch) Branch {
	return branch{
		branch: gitBranch,
	}
}
