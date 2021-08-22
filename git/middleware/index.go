package middleware

import git "github.com/libgit2/git2go/v31"

type Index interface {
	WriteTree() (*git.Oid, error)
}

type index struct {
	idx *git.Index
}

func (i index) WriteTree() (*git.Oid, error) {
	return i.idx.WriteTree()
}

func NewIndex(idx *git.Index) Index {
	return index{idx: idx}
}
