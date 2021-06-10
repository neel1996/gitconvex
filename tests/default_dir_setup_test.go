package tests

import (
	"github.com/neel1996/gitconvex/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultDirSetup_ShouldReturnDefaultPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		exec, _ := os.Executable()
		expectedPath := filepath.Dir(exec)

		execPath, err := utils.DefaultDirSetup()

		assert.Nil(t, err)
		assert.Equal(t, expectedPath, execPath)
	} else {
		expectedPath := "/usr/local/gitconvex"
		defaultPath, err := utils.DefaultDirSetup()

		assert.Nil(t, err)
		assert.Equal(t, expectedPath, defaultPath)
	}
}
