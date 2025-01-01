package unix

import (
	"os/exec"
	"strconv"
	"strings"
)

type Weekday int

const (
	Sunday    Weekday = 0
	Monday    Weekday = 1
	Tuesday   Weekday = 2
	Wednesday Weekday = 3
	Thursday  Weekday = 4
	Friday    Weekday = 5
	Saturday  Weekday = 6
)

// NewCronJob creates a new cron job with the specified command.
func NewCronJob(command string) CronJob {
	cron := new(cronDriver)
	cron.command = command
	cron.minute = "*"
	cron.hour = "*"
	cron.day = "*"
	cron.month = "*"
	cron.weekday = "*"
	return cron
}

// SetCronTZ sets the timezone of the cron daemon to the specified timezone.
func SetCronTZ(tz string) error {
	if lines, err := crons(); err != nil {
		return err
	} else {
		var result strings.Builder
		result.WriteString("TZ=" + tz + "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "TZ=") {
				result.WriteString(line + "\n")
			}
		}
		cmd := `echo "` + result.String() + `" | crontab -`
		return exec.Command("sudo", "bash", "-c", cmd).Run()
	}
}

// CronJob represents a cron job.
// .---------------- minute (0 - 59)
// |  .------------- hour (0 - 23)
// |  |  .---------- day of month (1 - 31)
// |  |  |  .------- month (1 - 12) OR jan,feb,mar,apr ...
// |  |  |  |  .---- day of week (0 - 6) (Sunday=0 or 7) OR sun,mon,tue,wed,thu,fri,sat
// |  |  |  |  |
// m h dom mon dow usercommand
type CronJob interface {
	// AtReboot schedules the cron job to run at reboot.
	AtReboot() CronJob
	// Yearly schedules the cron job to run every year.
	Yearly() CronJob
	// Monthly schedules the cron job to run every month.
	Monthly() CronJob
	// Weekly schedules the cron job to run every week.
	Weekly() CronJob
	// Daily schedules the cron job to run every day.
	Daily() CronJob
	// EveryXHours schedules the cron job to run every specified number of hours.
	EveryXHours(hours int) CronJob
	// EveryXMinutes schedules the cron job to run every specified number of minutes.
	EveryXMinutes(minutes int) CronJob
	// SetMinute sets the minute of the cron job.
	SetMinute(minute int) CronJob
	// SetHour sets the hour of the cron job.
	SetHour(hour int) CronJob
	// SetDayOfMonth sets the day of the month of the cron job.
	SetDayOfMonth(day int) CronJob
	// SetMonth sets the month of the cron job.
	SetMonth(month int) CronJob
	// SetDayOfWeek sets the day of the week of the cron job.
	SetDayOfWeek(day Weekday) CronJob
	// Command sets the command to be executed by the cron job.
	Command(command string) CronJob
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
	reboot  bool
	minute  string
	hour    string
	day     string
	month   string
	weekday string
	command string
}

func (cron *cronDriver) set(minute, hour, day, mon, wd string) CronJob {
	cron.minute = minute
	cron.hour = hour
	cron.day = day
	cron.month = mon
	cron.weekday = wd
	return cron
}

func (cron *cronDriver) AtReboot() CronJob {
	cron.reboot = true
	return cron
}

func (cron *cronDriver) Yearly() CronJob {
	return cron.set("0", "0", "1", "1", "*")
}

func (cron *cronDriver) Monthly() CronJob {
	return cron.set("0", "0", "1", "*", "*")
}

func (cron *cronDriver) Weekly() CronJob {
	return cron.set("0", "0", "*", "*", "0")
}

func (cron *cronDriver) Daily() CronJob {
	return cron.set("0", "0", "*", "*", "*")
}

func (cron *cronDriver) EveryXHours(hours int) CronJob {
	cron.hour = "*/" + strconv.Itoa(hours)
	return cron
}

func (cron *cronDriver) EveryXMinutes(minutes int) CronJob {
	cron.minute = "*/" + strconv.Itoa(minutes)
	return cron
}

func (cron *cronDriver) SetMinute(minute int) CronJob {
	if minute >= 0 && minute <= 59 {
		cron.minute = strconv.Itoa(minute)
	}
	return cron
}

func (cron *cronDriver) SetHour(hour int) CronJob {
	if hour >= 0 && hour <= 23 {
		cron.hour = strconv.Itoa(hour)
	}
	return cron
}

func (cron *cronDriver) SetDayOfMonth(day int) CronJob {
	if day >= 1 && day <= 31 {
		cron.day = strconv.Itoa(day)
	}
	return cron
}

func (cron *cronDriver) SetMonth(month int) CronJob {
	if month >= 1 && month <= 12 {
		cron.month = strconv.Itoa(month)
	}
	return cron
}

func (cron *cronDriver) SetDayOfWeek(day Weekday) CronJob {
	if day >= 0 && day <= 6 {
		cron.weekday = strconv.Itoa(int(day))
	}
	return cron
}

func (cron *cronDriver) Command(command string) CronJob {
	cron.command = command
	return cron
}

func (cron cronDriver) Compile() string {
	if cron.reboot {
		return "@reboot " + cron.command
	} else {
		return cron.minute + " " +
			cron.hour + " " +
			cron.day + " " +
			cron.month + " " +
			cron.weekday + " " +
			cron.command
	}
}

func (cron *cronDriver) Exists() (bool, error) {
	if lines, err := crons(); err != nil {
		return false, err
	} else {
		for _, line := range lines {
			if ok, cmd := cronCommand(line); ok && cmd == cron.command {
				return true, nil
			}
		}
		return false, nil
	}
}

func (cron *cronDriver) Install() (bool, error) {
	if lines, err := crons(); err != nil {
		return false, err
	} else {
		var result strings.Builder
		found := false
		for _, line := range lines {
			if ok, cmd := cronCommand(line); ok && cmd == cron.command {
				found = true
				result.WriteString(cron.Compile() + "\n")
			} else if strings.TrimSpace(line) != "" {
				result.WriteString(line + "\n")
			}
		}
		if !found {
			result.WriteString(cron.Compile() + "\n")
		}
		cmd := `echo "` + result.String() + `" | crontab -`
		return true, eOf(exec.Command("sudo", "bash", "-c", cmd).Run())
	}
}

func (cron *cronDriver) Uninstall() (bool, error) {
	if lines, err := crons(); err != nil {
		return false, err
	} else {
		var result strings.Builder
		for _, line := range lines {
			if ok, cmd := cronCommand(line); ok && cmd == cron.command {
				continue
			} else if strings.TrimSpace(line) != "" {
				result.WriteString(line + "\n")
			}
		}
		cmd := `echo "` + result.String() + `" | crontab -`
		return true, eOf(exec.Command("sudo", "bash", "-c", cmd).Run())
	}
}
