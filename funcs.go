package unix

import (
	"fmt"
	"os"
	"os/exec"
)

// RunAsSudo checks if the program is running with sudo privileges.
func RunAsSudo() bool {
	return os.Geteuid() == 0
}

// FileExists checks if a file exists and returns an error if it does not.
func FileExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// execE handle execute error for exit code
func execE(err error) error {
	if err == nil {
		return nil
	} else if exitErr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("%s", string(exitErr.Stderr))
	} else {
		return err
	}
}

// execVE handle execute error for command with output
func execVE[T any](v T, err error) (T, error) {
	if err == nil {
		return v, nil
	} else if exitErr, ok := err.(*exec.ExitError); ok {
		return v, fmt.Errorf("%s", string(exitErr.Stderr))
	} else {
		return v, err
	}
}
