package tests

import (
	"github.com/neel1996/gitconvex/git"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"testing"
)

var TestRepo string
var logger *logrus.Logger

func TestMain(m *testing.M) {
	logger = logrus.New()
	TestRepo = os.Getenv("GITCONVEX_TEST_REPO")
	if TestRepo == "" {
		logger.Error("No test repo passed. Set GITCONVEX_TEST_REPO to a git directory to run integration tests")
		os.Exit(1)
	}

	var testObject git.CloneInterface
	testObject = git.CloneStruct{
		RepoName:   "",
		RepoPath:   TestRepo,
		RepoURL:    "https://github.com/neel1996/starfleet.git",
		AuthOption: "noauth",
		SSHKeyPath: "",
		UserName:   "",
		Password:   "",
	}
	_, err := testObject.CloneRepo()
	if err != nil {
		logger.Error("Unable to clone the test repo")
		os.Exit(2)
	}

	logger.Info("Initiating integration tests")
	m.Run()
	tearDownTests()
}

func tearDownTests() {
	removeFlag, _ := strconv.ParseBool(os.Getenv("GITCONVEX_DELETE_TEST_REPO"))
	if !removeFlag {
		logger.Error("Test repo will be retained. If you use the same test repo again then the tests will fail during cloning stage")
		return
	}

	err := os.RemoveAll(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		logger.Error(err)
		return
	}
}
