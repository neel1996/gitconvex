package middleware

import git "github.com/libgit2/git2go/v31"

type Repository interface {
	Remotes() Remotes
	Walk() (RevWalk, error)
	Head() (Reference, error)
	LookupCommit(oid *git.Oid) (*git.Commit, error)
	LookupCommitV2(oid *git.Oid) (Commit, error)
	DefaultSignature() (*git.Signature, error)
	LookupTree(id *git.Oid) (*git.Tree, error)
	Index() (Index, error)
	CreateCommit(s string, signature *git.Signature, signature2 *git.Signature, message string, tree *git.Tree, parents ...*git.Commit) (*git.Oid, error)
	DiffTreeToTree(tree *git.Tree, tree2 *git.Tree, options *git.DiffOptions) (*git.Diff, error)
	CreateBranch(string, *git.Commit, bool) (*git.Branch, error)
	LookupBranch(branchName string, branchType git.BranchType) (Branch, error)
	CheckoutTree(tree *git.Tree, c *git.CheckoutOptions) error
	SetHead(name string) error
	GetGitRepository() *git.Repository
	NewBranchIterator(branchType git.BranchType) (BranchIterator, error)
}

type repository struct {
	repo *git.Repository
}

func (r repository) NewBranchIterator(branchType git.BranchType) (BranchIterator, error) {
	itr, err := r.repo.NewBranchIterator(branchType)
	if err != nil {
		return nil, err
	}

	return NewBranchIterator(itr), nil
}

func (r repository) SetHead(name string) error {
	return r.repo.SetHead(name)
}

func (r repository) CheckoutTree(tree *git.Tree, c *git.CheckoutOptions) error {
	return r.repo.CheckoutTree(tree, c)
}

func (r repository) Remotes() Remotes {
	return NewRemotes(r.repo.Remotes)
}

func (r repository) CreateCommit(s string, signature *git.Signature, signature2 *git.Signature,
	message string, tree *git.Tree, parents ...*git.Commit) (*git.Oid, error) {
	return r.repo.CreateCommit(s, signature, signature2, message, tree, parents...)
}

func (r repository) CreateBranch(branchName string, commit *git.Commit, force bool) (*git.Branch, error) {
	return r.repo.CreateBranch(branchName, commit, force)
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

func (r repository) LookupCommitV2(oid *git.Oid) (Commit, error) {
	gitCommit, err := r.repo.LookupCommit(oid)
	if err != nil {
		return nil, err
	}
	return NewCommit(gitCommit), nil
}

func (r repository) LookupBranch(branchName string, branchType git.BranchType) (Branch, error) {
	gitBranch, err := r.repo.LookupBranch(branchName, branchType)
	if err != nil {
		return nil, err
	}
	return NewBranch(gitBranch), nil
}

func (r repository) LookupTree(id *git.Oid) (*git.Tree, error) {
	return r.repo.LookupTree(id)
}

func (r repository) Walk() (RevWalk, error) {
	walk, err := r.repo.Walk()

	return NewRevWalk(walk), err
}

func (r repository) DefaultSignature() (*git.Signature, error) {
	return r.repo.DefaultSignature()
}

func (r repository) DiffTreeToTree(tree *git.Tree, tree2 *git.Tree, options *git.DiffOptions) (*git.Diff, error) {
	return r.repo.DiffTreeToTree(tree, tree2, options)
}

func (r repository) GetGitRepository() *git.Repository {
	return r.repo
}

func NewRepository(repo *git.Repository) Repository {
	return repository{repo: repo}
}
