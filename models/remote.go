package models

import "github.com/neel1996/gitconvex/git/remote"

type Remote struct {
	AddRemote     remote.Add
	DeleteRemote  remote.Delete
	EditRemote    remote.Edit
	ListRemote    remote.List
	RemoteName    remote.Name
	ListRemoteUrl remote.ListRemoteUrl
}
