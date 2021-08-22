package commit

import (
	"fmt"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"strings"
)

type LatestMessage interface {
	Get() string
}

type latestMessage struct {
	repo middleware.Repository
}

func (l latestMessage) Get() string {
	head, err := l.repo.Head()
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return ""
	}

	commit := head.Target()

	lookupCommit, lookupError := l.repo.LookupCommit(commit)
	if lookupError != nil {
		logger.Log("Unable to lookup commit : "+lookupError.Error(), global.StatusError)
		return ""
	}

	commitMessage := lookupCommit.Message()
	if commitMessage != "" {
		return l.truncateMultilineMessage(commitMessage)
	}

	return commitMessage
}

func (l latestMessage) truncateMultilineMessage(commitMessage string) string {
	if strings.Contains(commitMessage, "\n") {
		logger.Log("Truncating multiline commit message", global.StatusInfo)
		commitMessage = strings.Split(commitMessage, "\n")[0]
	}

	logger.Log(fmt.Sprintf("HEAD commit message is : %v", commitMessage), global.StatusInfo)
	return commitMessage
}

func NewLatestMessage(repo middleware.Repository) LatestMessage {
	return latestMessage{repo: repo}
}
