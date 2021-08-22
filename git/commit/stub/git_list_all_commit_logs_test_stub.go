package stub

import (
	"errors"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
)

type revWalkStub struct {
	walk                     *git2go.RevWalk
	shouldIterateReturnError bool
}

func (w *revWalkStub) Push(id *git2go.Oid) error {
	return nil
}

func (w *revWalkStub) PushHead() error {
	return nil
}

func (w *revWalkStub) Iterate(iterator git2go.RevWalkIterator) error {
	if w.shouldIterateReturnError {
		return errors.New("ITERATE_ERROR")
	}
	return nil
}

func NewRevWalkStub(shouldIterateReturnError bool) middleware.RevWalk {
	return &revWalkStub{walk: &git2go.RevWalk{}, shouldIterateReturnError: shouldIterateReturnError}
}
