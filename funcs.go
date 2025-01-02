package unix

import (
	"os"
	"strings"
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

// QuickReplace replaces all {instances} of the search string with the replacement string.
func QuickReplace(content string, replacements ...string) string {
	return strings.NewReplacer(replacements...).Replace(content)
}
