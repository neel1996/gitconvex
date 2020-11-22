package git

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/neel1996/gitconvex-server/global"
	"go/types"
	"strings"
)

type Branch struct {
	CurrentBranch string
	BranchList    []*string
	AllBranchList []*string
}

// GetBranchList fetches all the branches from the target repository
// The result will be returned as a struct with the current branch and all the available branches
func GetBranchList(repo *git.Repository, branchChan chan Branch) {
	var (
		branches      []*string
		allBranchList []*string
	)
	var currentBranch string

	logger := global.Logger{}

	if repo != nil {
		head, headErr := repo.Head()

		if headErr != nil {
			logger.Log(fmt.Sprintf("Repo head is invalid -> %s", headErr.Error()), global.StatusError)
			nilMsg := []string{"No Branches available", "Repo HEAD is nil", "No Branches Present"}
			branchChan <- Branch{
				BranchList:    []*string{&nilMsg[0]},
				CurrentBranch: nilMsg[1],
				AllBranchList: []*string{&nilMsg[2]},
			}
		} else {
			currentBranch = head.Name().String()
			splitCurrentBranch := strings.Split(currentBranch, "/")
			currentBranch = splitCurrentBranch[len(splitCurrentBranch)-1]

			ref, _ := repo.References()

			if ref != nil {
				_ = ref.ForEach(func(reference *plumbing.Reference) error {
					var (
						refNamePtr *string
					)
					if ref != nil {
						if reference.Name().String() != "HEAD" && strings.Contains(reference.Name().String(), "refs/") {
							refNameSplit := strings.Split(reference.Name().String(), "refs/")
							if len(refNameSplit) == 2 && strings.Contains(refNameSplit[1], "/") && !strings.Contains(refNameSplit[1], "remotes/"+git.DefaultRemoteName+"/HEAD") {
								logger.Log(fmt.Sprintf("Available Branch : %v", refNameSplit[1]), global.StatusInfo)
								if strings.Contains(refNameSplit[1], "heads/") {
									headBranch := strings.Split(refNameSplit[1], "heads/")[1]
									refNamePtr = &headBranch
								} else {
									refNamePtr = &refNameSplit[1]
								}
								allBranchList = append(allBranchList, refNamePtr)
							}
						}
					}
					return nil
				})
			} else {
				logger.Log("No references found!", global.StatusError)
			}

			bIter, _ := repo.Branches()

			if bIter != nil {
				_ = bIter.ForEach(func(reference *plumbing.Reference) error {
					if reference != nil {
						localBranch := reference.String()
						splitBranch := strings.Split(localBranch, "/")
						localBranch = splitBranch[len(splitBranch)-1]

						branches = append(branches, &localBranch)
						return nil
					} else {
						return types.Error{Msg: "Empty reference"}
					}
				})
				bIter.Close()
			} else {
				logger.Log("Nil branch reference found", global.StatusError)
			}
		}

		branchChan <- Branch{
			BranchList:    branches,
			CurrentBranch: currentBranch,
			AllBranchList: allBranchList,
		}
	}

	close(branchChan)
}
