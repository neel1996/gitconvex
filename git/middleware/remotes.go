package middleware

import git "github.com/libgit2/git2go/v31"

type Remotes interface {
	Create(string, string) (*git.Remote, error)
	Get() git.RemoteCollection
	Delete(string) error
	List() ([]string, error)
	SetUrl(string, string) error
	Lookup(string) (*git.Remote, error)
}

type remotes struct {
	git.RemoteCollection
}

func (r remotes) Create(name string, url string) (*git.Remote, error) {
	return r.RemoteCollection.Create(name, url)
}

func (r remotes) Get() git.RemoteCollection {
	return r.RemoteCollection
}

func (r remotes) Delete(name string) error {
	return r.RemoteCollection.Delete(name)
}

func (r remotes) List() ([]string, error) {
	return r.RemoteCollection.List()
}

func (r remotes) SetUrl(name string, url string) error {
	return r.RemoteCollection.SetUrl(name, url)
}

func (r remotes) Lookup(name string) (*git.Remote, error) {
	return r.RemoteCollection.Lookup(name)
}

func NewRemotes(remoteCollection git.RemoteCollection) Remotes {
	return remotes{
		RemoteCollection: remoteCollection,
	}
}
