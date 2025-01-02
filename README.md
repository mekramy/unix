# Unix Package Documentation

## Overview

The `unix` package provides utilities for managing system services, cron jobs, and Nginx reverse proxies on Unix-based systems. It includes functionalities to check file existence, manage systemd services, configure Nginx server blocks, and handle cron jobs.

## Installation

To install the package, use the following command:

```sh
go get github.com/mekramy/unix
```

## Functions

### RunAsSudo

```go
func RunAsSudo() bool
```

Checks if the program is running with sudo privileges.

### FileExists

```go
func FileExists(filePath string) (bool, error)
```

Checks if a file exists and returns an error if it does not.

### QuickReplace

```go
func QuickReplace(content string, replacements ...string) string
```

replaces all {instances} of the search string with the replacement string.

## Systemd Service Management

### NewSystemdService

```go
func NewSystemdService(name, root, command string) SystemdService
```

Creates a new systemd service.

### SystemdService Interface

- `Name(name string) SystemdService`: Sets the name of the service.
- `Root(dir string) SystemdService`: Sets the root path of the service.
- `Command(command string) SystemdService`: Sets the command of the service.
- `Template(engine TemplateEngine) SystemdService`: Sets the template for the service.
- `Exists() bool`: Checks if the service exists.
- `Enabled() (bool, error)`: Checks if the service exists and is enabled on startup.
- `Install(override bool) (bool, error)`: Installs the service.
- `Uninstall() error`: Uninstalls the service.

## Nginx Reverse Proxy Management

### NewNginxReverseProxy

```go
func NewNginxReverseProxy(name, port string) ServerBlock
```

Creates a new Nginx reverse proxy server block.

### ServerBlock Interface

- `Name(name string) ServerBlock`: Sets the name of the site.
- `Port(port string) ServerBlock`: Sets the port for the site.
- `Domains(domains ...string) ServerBlock`: Sets the domains for the site.
- `Template(engine TemplateEngine) ServerBlock`: Sets the template for the site.
- `Disable() error`: Disables the site manually.
- `Enable() error`: Enables the site manually.
- `Exists() (bool, error)`: Checks if the site exists.
- `Enabled() (bool, error)`: Checks if the site exists and is enabled.
- `Install(override bool) (bool, error)`: Installs the site.
- `Uninstall() error`: Uninstalls the site.

## Cron Job Management

### NewCronJob

```go
func NewCronJob(command string) CronJob
```

Creates a new cron job.

### CronJob Interface

- `SetTz(hour int, min int) CronJob`: sets the timezone of the cron job.
- `AtReboot() CronJob`: schedules the cron job to run at reboot.
- `Yearly() CronJob`: schedules the cron job to run every year.
- `Monthly() CronJob`: schedules the cron job to run every month.
- `Weekly(wd Weekday) CronJob`: schedules the cron job to run every week.
- `Daily() CronJob`: schedules the cron job to run every day.
- `EveryXHours(hours int) CronJob`: schedules the cron job to run every specified number of hours.
- `EveryXMinutes(minutes int) CronJob`: schedules the cron job to run every specified number of minutes.
- `SetMinute(minute int) CronJob`: sets the minute of the cron job.
- `SetHour(hour int) CronJob`: sets the hour of the cron job.
- `SetDayOfMonth(day int) CronJob`: sets the day of the month of the cron job.
- `SetMonth(month int) CronJob`: sets the month of the cron job.
- `SetDayOfWeek(day Weekday) CronJob`: sets the day of the week of the cron job.
- `Command(command string) CronJob`: sets the command to be executed by the cron job.
- `Compile() string`: compiles the cron job into a cron expression string.
- `Exists() (bool, error)`: checks if the cron job already exists.
- `Install() (bool, error)`: installs the cron job. returns false if cronjob exists.
- `Uninstall() error`: uninstalls the cron job.

## Formatter Utility

### PrintF

```go
func PrintF(format string, args ...any)
```

The `PrintF` print stylized text to console. The pattern string can contain tags followed by content tokens for various styles and colors. You can escape tokens with `\@`.

#### Supported Styling Patterns

- `B`: BOLD
- `U`: UNDERLINE
- `S`: STRIKE
- `I`: ITALIC

#### Supported Color Patterns

- `r`: RED
- `g`: GREEN
- `y`: YELLOW
- `b`: BLUE
- `p`: PURPLE
- `c`: CYAN
- `m`: GRAY
- `w`: WHITE

#### Example Usage

```go
PrintF("@Bg{Bold Green Text} and @r{Red %s}\n", "message")
```

#### Arguments

- `format`: The string containing the standard Go `fmt` format with styled tokens.
- `args`: The arguments to be passed into the format string.

## License

This project is licensed under the MIT License.
