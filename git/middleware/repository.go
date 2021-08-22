package middleware

import git "github.com/libgit2/git2go/v31"

type Repository interface {
	Remotes() Remotes
	Walk() (RevWalk, error)
	Head() (Reference, error)
	LookupCommit(oid *git.Oid) (*git.Commit, error)
	DefaultSignature() (*git.Signature, error)
	LookupTree(id *git.Oid) (*git.Tree, error)
	Index() (Index, error)
	CreateCommit(s string, signature *git.Signature, signature2 *git.Signature, message string, tree *git.Tree, parents ...*git.Commit) (*git.Oid, error)
	DiffTreeToTree(tree *git.Tree, tree2 *git.Tree, options *git.DiffOptions) (*git.Diff, error)
}

type repository struct {
	repo *git.Repository
}

func (r repository) Remotes() Remotes {
	return NewRemotes(r.repo.Remotes)
}

func (r repository) CreateCommit(s string, signature *git.Signature, signature2 *git.Signature,
	message string, tree *git.Tree, parents ...*git.Commit) (*git.Oid, error) {
	return r.repo.CreateCommit(s, signature, signature2, message, tree, parents...)
}

func (r repository) Head() (Reference, error) {
	head, err := r.repo.Head()
	if err != nil {
		return nil, err
	}
	return NewReference(head), nil
}

func (r repository) Index() (Index, error) {
	idx, err := r.repo.Index()
	if err != nil {
		return idx, err
	}
	return NewIndex(idx), nil
}

func (r repository) LookupCommit(oid *git.Oid) (*git.Commit, error) {
	return r.repo.LookupCommit(oid)
}

func (r repository) Walk() (RevWalk, error) {
	walk, err := r.repo.Walk()

	return NewRevWalk(walk), err
}

func (r repository) DefaultSignature() (*git.Signature, error) {
	return r.repo.DefaultSignature()
}

func (r repository) LookupTree(id *git.Oid) (*git.Tree, error) {
	return r.repo.LookupTree(id)
}

func (r repository) DiffTreeToTree(tree *git.Tree, tree2 *git.Tree, options *git.DiffOptions) (*git.Diff, error) {
	return r.repo.DiffTreeToTree(tree, tree2, options)
}

func NewRepository(repo *git.Repository) Repository {
	return repository{repo: repo}
}
