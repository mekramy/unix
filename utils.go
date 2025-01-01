package unix

import (
	"fmt"
	"os/exec"
	"strings"
)

// eOf handle execute error for exit code
func eOf(err error) error {
	if err == nil {
		return nil
	} else if exitErr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("%s", string(exitErr.Stderr))
	} else {
		return err
	}
}

// evOf handle execute error for command with output
func evOf[T any](v T, err error) (T, error) {
	return v, eOf(err)
}

// crons get all cron jobs
func crons() ([]string, error) {
	out, err := evOf(exec.Command("sudo", "crontab", "-l").Output())
	if err != nil {
		return nil, err
	}
	return strings.Split(string(out), "\n"), nil
}

// cronCommand extracts the command from a cron expression.
func cronCommand(cronExpr string) (bool, string) {
	parts := strings.Fields(cronExpr)
	if len(parts) < 6 {
		return false, ""
	}
	return true, strings.Join(parts[5:], " ")
}
