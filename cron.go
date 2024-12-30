package unix

import (
	"os/exec"
	"strconv"
	"strings"
)

func NewCronJob(command string) CronJob {
	cron := new(cronDriver)
	cron.command = command
	return cron
}

type CronJob interface {
	// Command sets the command to be executed by the cron job.
	Command(command string) CronJob
	// UTCTz sets the timezone of the cron job to UTC.
	UTCTz() CronJob
	// TehranTz sets the timezone of the cron job to Asia/Tehran.
	TehranTz() CronJob
	// Timezone sets the timezone of the cron job to the specified timezone.
	Timezone(tz string) CronJob
	// Hourly schedules the cron job to run every hour.
	Hourly() CronJob
	// Daily schedules the cron job to run every day.
	Daily() CronJob
	// Weekly schedules the cron job to run every week.
	Weekly() CronJob
	// Monthly schedules the cron job to run every month.
	Monthly() CronJob
	// EveryXHours schedules the cron job to run every specified number of hours.
	EveryXHours(hours int) CronJob
	// EveryXMinutes schedules the cron job to run every specified number of minutes.
	EveryXMinutes(minutes int) CronJob
	// Compile compiles the cron job into a cron expression string.
	Compile() string
	// Exists checks if the cron job already exists.
	Exists() (bool, error)
	// Install installs the cron job. returns false if cronjob exists.
	Install() (bool, error)
	// Uninstall uninstalls the cron job. returns false if cronjob does not exist.
	Uninstall() (bool, error)
}

type cronDriver struct {
	command  string
	timezone string
	interval string
}

func (cron *cronDriver) Command(command string) CronJob {
	cron.command = command
	return cron
}

func (cron *cronDriver) UTCTz() CronJob {
	cron.timezone = ""
	return cron
}

func (cron *cronDriver) TehranTz() CronJob {
	cron.timezone = "Asia/Tehran"
	return cron
}

func (cron *cronDriver) Timezone(tz string) CronJob {
	cron.timezone = tz
	return cron
}

func (cron *cronDriver) Hourly() CronJob {
	cron.interval = "0 * * * *"
	return cron
}

func (cron *cronDriver) Daily() CronJob {
	cron.interval = "0 0 * * *"
	return cron
}

func (cron *cronDriver) Weekly() CronJob {
	if cron.timezone == "Asia/Tehran" {
		cron.interval = "0 0 * * 5"
	} else {
		cron.interval = "0 0 * * 0"
	}
	return cron
}

func (cron *cronDriver) Monthly() CronJob {
	cron.interval = "0 * * * *"
	return cron
}

func (cron *cronDriver) EveryXHours(hours int) CronJob {
	cron.interval = "0 */" + strconv.Itoa(hours) + " * * *"
	return cron
}

func (cron *cronDriver) EveryXMinutes(minutes int) CronJob {
	cron.interval = "*/" + strconv.Itoa(minutes) + " * * * *"
	return cron
}

func (cron *cronDriver) Compile() string {
	if cron.timezone == "" {
		return cron.interval + " " + cron.command
	} else {
		return cron.interval + " " + cron.command + " TZ=" + cron.timezone
	}
}

func (cron *cronDriver) Exists() (bool, error) {
	if out, err := execVE(exec.Command("sudo", "crontab", "-l").Output()); err != nil {
		return false, err
	} else {
		return strings.Contains(string(out), cron.Compile()), nil
	}
}

func (cron *cronDriver) Install() (bool, error) {
	if exists, err := cron.Exists(); err != nil {
		return false, err
	} else if exists {
		return false, nil
	} else {
		cmd := `(crontab -l ; echo "` + cron.Compile() + `") | crontab -`
		return true, execE(exec.Command("sudo", "bash", "-c", cmd).Run())
	}
}

func (cron *cronDriver) Uninstall() (bool, error) {
	if exists, err := cron.Exists(); err != nil {
		return false, err
	} else if !exists {
		return false, nil
	} else {
		if out, err := execVE(exec.Command("sudo", "crontab", "-l").Output()); err != nil {
			return false, err
		} else {
			lines := strings.Split(string(out), "\n")
			newCron := ""
			for _, line := range lines {
				if line != "" && !strings.Contains(line, cron.Compile()) {
					newCron += line + "\n"
				}
			}
			cmd := `echo "` + newCron + `" | crontab -`
			return true, execE(exec.Command("sudo", "bash", "-c", cmd).Run())
		}
	}
}
