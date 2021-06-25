package remote

import "strings"

type Host interface {
	GetRemoteHostForUrl() string
}

type host struct {
	remoteUrl string
}

func (h host) GetRemoteHostForUrl() string {
	var remoteHostNames = []string{"github", "gitlab", "bitbucket", "azure", "codecommit"}

	for _, hostName := range remoteHostNames {
		if strings.Contains(h.remoteUrl, hostName) {
			return hostName
		}
	}

	return ""
}

func NewRemoteHost(remoteUrl string) Host {
	return host{remoteUrl: remoteUrl}
}
