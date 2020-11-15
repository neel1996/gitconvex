package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/nleeper/goment"
	"strconv"
	"strings"
	"time"
)

// commitOrganizer collects and organizes the commit related information in the commit GitCommits struct

func commitOrganizer(commits []object.Commit) []*model.GitCommits {
	logger := global.Logger{}
	var commitList []*model.GitCommits
	for _, commit := range commits {
		if !commit.Hash.IsZero() {
			commitHash := commit.Hash.String()
			commitAuthor := strings.Split(commit.Author.String(), " ")[0]
			commitMessage := strings.Split(commit.Message, "\n")[0]
			commitFilesItr, err := commit.Files()
			commitFileCount := 0
			commitDate := ""
			commitRelativeTime := ""

			logger.Log(fmt.Sprintf("Fetching commit details for -> %s", commitHash), global.StatusInfo)

			var prevTree *object.Tree
			prevCommit, parentErr := commit.Parents().Next()
			currentTree, _ := commit.Tree()

			if parentErr != nil {
				commitFileCount = 0
			} else {
				prevTree, _ = prevCommit.Tree()
				diff, _ := currentTree.Diff(prevTree)
				commitFileCount = diff.Len()
			}

			// Logic to extract commit timestamp from commit string
			// go-git commit iterator does not provide an option to extract the timestamp directly

			for _, cString := range strings.Split(commit.String(), "\n") {
				if strings.Contains(cString, "Date:") {
					str := strings.Split(cString, "Date:")[1]
					tempDate := strings.TrimSpace(str)

					if strings.Contains(tempDate, "+") {
						tempDate = strings.TrimSpace(strings.Split(tempDate, "+")[0])
					} else if strings.Contains(tempDate, "-") {
						tempDate = strings.TrimSpace(strings.Split(tempDate, "-")[0])
					}

					cTime, convErr := time.Parse(time.ANSIC, tempDate)
					if convErr != nil {
						logger.Log(convErr.Error(), global.StatusError)
					} else {
						commitDate = cTime.String()
						gTime, gTimeErr := goment.New(cTime)
						if gTimeErr != nil {
							logger.Log(gTimeErr.Error(), global.StatusError)
						} else {
							commitRelativeTime = gTime.FromNow()

							// Conditional logic to find time diff to bypass goment bug
							if strings.Contains(commitRelativeTime, "in") {
								aTime := time.Now().String()

								a, _ := time.Parse("2006-01-02 15:04:05", aTime[:19])
								b, _ := time.Parse("2006-01-02 15:04:05", cTime.String()[:19])
								diff := a.Sub(b)

								h := diff.Hours()
								m := diff.Minutes()
								s := diff.Seconds()

								if h != float64(0) {
									hStr := strconv.Itoa(int(h))
									if hStr == "1" {
										commitRelativeTime = hStr + " hour ago"
									} else {
										commitRelativeTime = hStr + " hours ago"
									}
								} else {
									if m != float64(0) {
										mStr := strconv.Itoa(int(m))
										commitRelativeTime = mStr + " minutes ago"
									} else {
										if s != float64(0) {
											sStr := strconv.Itoa(int(s))
											commitRelativeTime = sStr + " seconds ago"
										} else {
											commitRelativeTime = "recent commit"

										}
									}
								}
							}
						}
					}
				}
			}

			if err != nil {
				logger.Log(err.Error(), global.StatusError)
			} else {
				_ = commitFilesItr.ForEach(func(file *object.File) error {
					return nil
				})
			}

			commitList = append(commitList, &model.GitCommits{
				Hash:               &commitHash,
				Author:             &commitAuthor,
				CommitTime:         &commitDate,
				CommitMessage:      &commitMessage,
				CommitRelativeTime: &commitRelativeTime,
				CommitFilesCount:   &commitFileCount,
			})
		}
	}

	return commitList
}

// CommitLogs fetches the structured commit logs list for the target repo
// The skipCount limit is set to limit the number of commit logs returned per invokation

func CommitLogs(repo *git.Repository, skipCount int) *model.GitCommitLogResults {
	var commitLogs []object.Commit

	allCommitChan := make(chan AllCommitData)
	go AllCommits(repo, allCommitChan)
	acc := <-allCommitChan
	totalCommits := acc.TotalCommits

	head, _ := repo.Head()
	commitItr, commitErr := repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Order: git.LogOrderCommitterTime,
		All:   false,
	})

	if commitErr == nil {
		_ = commitItr.ForEach(func(commit *object.Commit) error {
			commitLogs = append(commitLogs, *commit)
			return nil
		})
	}

	if len(commitLogs) == 0 {
		return &model.GitCommitLogResults{
			TotalCommits: &totalCommits,
			Commits:      nil,
		}
	}

	if len(commitLogs) <= 10 {
		refinedCommits := commitOrganizer(commitLogs)
		return &model.GitCommitLogResults{
			TotalCommits: &totalCommits,
			Commits:      refinedCommits,
		}
	} else {
		var commitSlice []object.Commit

		commitLimit := skipCount + 10
		if commitLimit > len(commitLogs) {
			commitLimit = skipCount
			commitSlice = commitLogs[skipCount:]
		} else {
			commitSlice = commitLogs[skipCount:commitLimit]
		}
		refinedCommits := commitOrganizer(commitSlice)
		return &model.GitCommitLogResults{
			TotalCommits: &totalCommits,
			Commits:      refinedCommits,
		}
	}
}
