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

- `Command(command string) CronJob`: Sets the command to be executed by the cron job.
- `UTCTz() CronJob`: Sets the timezone of the cron job to UTC.
- `TehranTz() CronJob`: Sets the timezone of the cron job to Asia/Tehran.
- `Timezone(tz string) CronJob`: Sets the timezone of the cron job to the specified timezone.
- `Hourly() CronJob`: Schedules the cron job to run every hour.
- `Daily() CronJob`: Schedules the cron job to run every day.
- `Weekly() CronJob`: Schedules the cron job to run every week.
- `Monthly() CronJob`: Schedules the cron job to run every month.
- `EveryXHours(hours int) CronJob`: Schedules the cron job to run every specified number of hours.
- `EveryXMinutes(minutes int) CronJob`: Schedules the cron job to run every specified number of minutes.
- `Compile() string`: Compiles the cron job into a cron expression string.
- `Exists() (bool, error)`: Checks if the cron job already exists.
- `Install() (bool, error)`: Installs the cron job.
- `Uninstall() (bool, error)`: Uninstalls the cron job.

## Formatter Functions

### PrintF

```go
func PrintF(pattern string, args ...any)
```

Prints a formatted string according to the specified pattern. The pattern can include styling placeholders that will be replaced with corresponding ANSI escape codes.

#### Styling Patterns

- `{R}`, `@`: RESET
- `{B}`, `@B`: BOLD
- `{U}`, `@U`: UNDERLINE
- `{S}`, `@S`: STRIKE
- `{I}`, `@I`: ITALIC
- `{r}`, `@r`: RED
- `{g}`, `@g`: GREEN
- `{y}`, `@y`: YELLOW
- `{b}`, `@b`: BLUE
- `{p}`, `@p`: PURPLE
- `{c}`, `@c`: CYAN
- `{m}`, `@m`: GRAY
- `{w}`, `@w`: WHITE

#### Example Usage

```go
Formatter("{B}Bold Text{R} and {r}Red Text{R}\n")
Formatter("{g}Green Text{R} with arguments: %d, %s\n", 42, "example")
```

## License

This project is licensed under the MIT License.
