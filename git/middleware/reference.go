package middleware

import git "github.com/libgit2/git2go/v31"

type Reference interface {
	GetGitReference() *git.Reference
	Target() *git.Oid
	SetTarget(id *git.Oid, message string) (*git.Reference, error)
	Name() string
}

type reference struct {
	ref *git.Reference
}

func (r reference) Name() string {
	return r.ref.Name()
}

func (r reference) Target() *git.Oid {
	return r.ref.Target()
}

func (r reference) SetTarget(id *git.Oid, message string) (*git.Reference, error) {
	return r.ref.SetTarget(id, message)
}

func (r reference) GetGitReference() *git.Reference {
	return r.ref
}

func NewReference(ref *git.Reference) Reference {
	return reference{ref: ref}
}
