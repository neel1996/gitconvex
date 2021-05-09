package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/neel1996/gitconvex-server/utils"
	"go/types"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type AddRepoInterface interface {
	AddRepo() *model.AddRepoParams
	repoDataFileWriter(repoId string, repoAddStatus chan string)
}

type AddRepoInputs struct {
	RepoName    string
	RepoPath    string
	CloneSwitch bool
	RepoURL     string
	InitSwitch  bool
	AuthOption  string
	UserName    string
	Password    string
	SSHKeyPath  string
}

type RepoData struct {
	Id         string `json:"id"`
	RepoName   string `json:"repoName"`
	RepoPath   string `json:"repoPath"`
	TimeStamp  string `json:"timestamp"`
	AuthOption string `json:"authOption"`
	SSHKeyPath string `json:"sshKeyPath"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
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
		logger.Log(err.Error(), global.StatusError)
		return types.Error{Msg: err.Error()}
	}
	if dirErr != nil {
		logger.Log(fmt.Sprintf("Error occurred creating database directory \n%v", dirErr), global.StatusError)
		return types.Error{Msg: dirErr.Error()}
	}
	logger.Log("New repo datastore created successfully", global.StatusInfo)
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
func (inputs AddRepoInputs) repoDataFileWriter(repoId string, repoAddStatus chan string) {
	rArray := make([]RepoData, 1)

	if inputs.Password != "" {
		utfBytes := bytes.NewBufferString(inputs.Password)
		hashedBytes, _ := bcrypt.GenerateFromPassword(utfBytes.Bytes(), bcrypt.MinCost)
		if hashedBytes != nil {
			inputs.Password = string(hashedBytes)
		}
	}

	if inputs.Password != "" && repoId != "" {
		var encryptObject utils.PasswordCipherInterface
		encryptObject = utils.PasswordCipherStruct{
			PlainPassword: inputs.Password,
			KeyString:     repoId,
		}
		inputs.Password = encryptObject.EncryptPassword()
	}

	rArray[0] = RepoData{
		Id:         repoId,
		RepoName:   inputs.RepoName,
		RepoPath:   inputs.RepoPath,
		TimeStamp:  time.Now().String(),
		AuthOption: inputs.AuthOption,
		SSHKeyPath: inputs.SSHKeyPath,
		UserName:   inputs.UserName,
		Password:   inputs.Password,
	}

	envConfig := *utils.EnvConfigFileReader()

	dbFile := envConfig.DataBaseFile
	logger.Log("Opening DB file present in env_config", global.StatusInfo)
	_, fileOpenErr := os.Open(dbFile)

	if fileOpenErr != nil {
		logger.Log(fmt.Sprintf("Error occurred while opening repo data JSON file \n%v", fileOpenErr), global.StatusError)

		createErr := repoDataCreator(dbFile)

		if createErr != nil {
			logger.Log(createErr.Error(), global.StatusError)
			panic(createErr)
		} else {
			if err := dataFileWriteHandler(dbFile, rArray); err != nil && err.Error() != "" {
				logger.Log(err.Error(), global.StatusError)
				repoAddStatus <- "failed"
			} else {
				repoAddStatus <- "success"
			}
		}
	} else {
		if err := dataFileWriteHandler(dbFile, rArray); err != nil && err.Error() != "" {
			logger.Log(err.Error(), global.StatusError)
			repoAddStatus <- "failed"
		} else {
			repoAddStatus <- "success"
		}
	}
}

// AddRepo function gets the repository details and includes a record to the gitconvex repo datastore file
// If initSwitch is 'true' then the git repo init function will be invoked to initialize a new repo
// If cloneSwitch is 'true' then the repo will be cloned to the file system using the repoURL field
func (inputs AddRepoInputs) AddRepo() *model.AddRepoParams {
	repoName := inputs.RepoName
	repoPath := inputs.RepoPath
	cloneSwitch := inputs.CloneSwitch
	repoURL := inputs.RepoURL
	initSwitch := inputs.InitSwitch
	authOption := inputs.AuthOption
	userName := inputs.UserName
	password := inputs.Password
	sshKeyPath := inputs.SSHKeyPath

	if cloneSwitch && len(repoURL) > 0 {
		currentOs := HealthCheckApi().Os
		if currentOs == "windows" {
			repoPath = repoPath + "\\" + repoName
		} else {
			repoPath = repoPath + "/" + repoName
		}

		inputs.RepoPath = repoPath

		var cloneObject git.CloneInterface
		cloneObject = git.CloneStruct{
			RepoName:   repoName,
			RepoPath:   repoPath,
			RepoURL:    repoURL,
			AuthOption: authOption,
			UserName:   userName,
			Password:   password,
			SSHKeyPath: sshKeyPath,
		}

		_, err := cloneObject.CloneRepo()
		if err != nil {
			logger.Log(fmt.Sprintf("%v", err), global.StatusError)
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
			logger.Log(fmt.Sprintf("%v", err), global.StatusError)
			return &model.AddRepoParams{
				RepoID:  "",
				Status:  "Failed",
				Message: err.Error(),
			}
		}
	}

	_, invalidRepoErr := git.RepoValidator(repoPath)

	if invalidRepoErr != nil && !initSwitch {
		logger.Log(fmt.Sprintf("The repo is not a valid git repo\n%v", invalidRepoErr), global.StatusError)

		return &model.AddRepoParams{
			RepoID:  "",
			Status:  "Failed",
			Message: invalidRepoErr.Error(),
		}
	}

	var repoIdChannel = make(chan string)
	go repoIdGenerator(repoIdChannel)
	repoId := <-repoIdChannel

	var repoAddStatusChannel = make(chan string)
	go inputs.repoDataFileWriter(repoId, repoAddStatusChannel)
	status := <-repoAddStatusChannel

	close(repoIdChannel)
	close(repoAddStatusChannel)

	if status == "success" {
		logger.Log("Repo entry added to the data store", global.StatusInfo)
		return &model.AddRepoParams{
			RepoID:  repoId,
			Status:  "Repo Added",
			Message: "The new repository has been added to Gitconvex",
		}
	} else {
		logger.Log("Failed to add new repo entry", global.StatusError)
		return &model.AddRepoParams{
			RepoID: "",
			Status: "Failed",
		}
	}

}
