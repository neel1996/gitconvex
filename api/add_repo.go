package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/neel1996/gitconvex-server/utils"
	"go/types"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type RepoData struct {
	Id        string `json:"id"`
	RepoName  string `json:"repoName"`
	RepoPath  string `json:"repoPath"`
	TimeStamp string `json:"timestamp"`
}

// localLogger logs messages to the global logger module

func localLogger(message string, status string) {
	logger := &global.Logger{Message: message}
	logger.Log(logger.Message, status)
}

// repoIdGenerator generates a unique ID for the newly added repo

func repoIdGenerator(c chan string) {
	newUUID, _ := uuid.NewUUID()
	repoId := strings.Split(newUUID.String(), "-")[0]
	c <- repoId
}

// repoDataCreator creates a new datastore directory and file if repo data does not exist

func repoDataCreator(dbFile string) error {
	sliceDbFile := strings.Split(dbFile, "/")
	dbDir := strings.Join(sliceDbFile[0:len(sliceDbFile)-1], "/")
	dirErr := os.MkdirAll(dbDir, 0755)
	_, err := os.Create(dbFile)

	if err != nil {
		localLogger(err.Error(), global.StatusError)
		return types.Error{Msg: err.Error()}
	}
	if dirErr != nil {
		localLogger(fmt.Sprintf("Error occurred creating database directory \n%v", dirErr), global.StatusError)
		return types.Error{Msg: dirErr.Error()}
	}
	localLogger("New repo datastore created successfully", global.StatusInfo)
	return nil
}

func dataFileWriteHandler(dbFile string, repoDataArray []RepoData) error {
	repoDataJSON, _ := json.Marshal(repoDataArray)
	osRead, _ := os.Open(dbFile)
	readFileStat, _ := ioutil.ReadAll(osRead)

	var existingData []RepoData

	if readFileStat != nil && json.Unmarshal(readFileStat, &existingData) == nil {
		appendData := append(existingData, repoDataArray[0])
		writeData, _ := json.MarshalIndent(appendData, "", " ")
		return ioutil.WriteFile(dbFile, writeData, 0755)
	} else {
		return ioutil.WriteFile(dbFile, repoDataJSON, 0755)
	}
}

// repoDataFileWriter writes the new repo details to the repo_datastore.json file

func repoDataFileWriter(repoId string, repoName string, repoPath string, repoAddStatus chan string) {
	rArray := make([]RepoData, 1)

	rArray[0] = RepoData{
		Id:        repoId,
		RepoName:  repoName,
		RepoPath:  repoPath,
		TimeStamp: time.Now().String(),
	}

	envConfig := *utils.EnvConfigFileReader()

	dbFile := envConfig.DataBaseFile
	localLogger("Opening DB file present in env_config", global.StatusInfo)
	_, fileOpenErr := os.Open(dbFile)

	if fileOpenErr != nil {
		localLogger(fmt.Sprintf("Error occurred while opening repo data JSON file \n%v", fileOpenErr), global.StatusError)

		createErr := repoDataCreator(dbFile)

		if createErr != nil {
			localLogger(createErr.Error(), global.StatusError)
			panic(createErr)
		} else {
			if err := dataFileWriteHandler(dbFile, rArray); err != nil && err.Error() != "" {
				localLogger(err.Error(), global.StatusError)
				repoAddStatus <- "failed"
			} else {
				repoAddStatus <- "success"
			}
		}
	} else {
		if err := dataFileWriteHandler(dbFile, rArray); err != nil && err.Error() != "" {
			localLogger(err.Error(), global.StatusError)
			repoAddStatus <- "failed"
		} else {
			repoAddStatus <- "success"
		}
	}
}

// AddRepo function gets the repository details and includes a record to the gitconvex repo datastore file
// If initSwitch is 'true' then the git repo init function will be invoked to initialize a new repo
// If cloneSwitch is 'true' then the repo will be cloned to the file system using the repoURL field

func AddRepo(inputs model.NewRepoInputs) *model.AddRepoParams {
	var repoIdChannel = make(chan string)

	repoName := inputs.RepoName
	repoPath := inputs.RepoPath
	cloneSwitch := inputs.CloneSwitch
	repoURL := inputs.RepoURL
	initSwitch := inputs.InitSwitch
	authOption := inputs.AuthOption
	userName := inputs.UserName
	password := inputs.Password

	if cloneSwitch && len(*repoURL) > 0 {
		repoPath = repoPath + "/" + repoName

		_, err := git.CloneHandler(repoPath, *repoURL, authOption, userName, password)
		if err != nil {
			localLogger(fmt.Sprintf("%v", err), global.StatusError)
			return &model.AddRepoParams{
				RepoID:  "",
				Status:  "Failed",
				Message: err.Error(),
			}
		}
	}
	if initSwitch {
		_, err := git.InitHandler(repoPath)
		if err != nil {
			localLogger(fmt.Sprintf("%v", err), global.StatusError)
			return &model.AddRepoParams{
				RepoID:  "",
				Status:  "Failed",
				Message: err.Error(),
			}
		}
	}

	_, invalidRepoErr := git.RepoValidator(repoPath)

	if invalidRepoErr != nil && !initSwitch {
		localLogger(fmt.Sprintf("The repo is not a valid git repo\n%v", invalidRepoErr), global.StatusError)

		return &model.AddRepoParams{
			RepoID:  "",
			Status:  "Failed",
			Message: invalidRepoErr.Error(),
		}
	}

	go repoIdGenerator(repoIdChannel)
	repoId := <-repoIdChannel

	var repoAddStatusChannel = make(chan string)
	go repoDataFileWriter(repoId, repoName, repoPath, repoAddStatusChannel)
	status := <-repoAddStatusChannel

	close(repoIdChannel)
	close(repoAddStatusChannel)

	if status == "success" {
		localLogger("Repo entry added to the data store", global.StatusInfo)
		return &model.AddRepoParams{
			RepoID:  repoId,
			Status:  "Repo Added",
			Message: "The new repository has been added to Gitconvex",
		}
	} else {
		localLogger("Failed to add new repo entry", global.StatusError)
		return &model.AddRepoParams{
			RepoID: "",
			Status: "Failed",
		}
	}

}
