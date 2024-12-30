package unix

import (
	"os"
	"os/exec"
	"strings"
)

func NewSystemdService(name, root, command string) SystemdService {
	service := new(systemdDriver)
	service.name = name
	service.root = root
	service.command = command
	service.template = NewEngine()
	service.template.SetTemplate(`
[Unit]
Description={name}
ConditionPathExists={root}
After=network.target

[Service]
Type=simple
User=root
Group=root
LimitNOFILE=1024

Restart=on-failure
RestartSec=10

WorkingDirectory={root}
ExecStart=/usr/bin/sudo {root}/{command}

PermissionsStartOnly=true
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier={name}

[Install]
WantedBy=multi-user.target
	`)
	return service
}

type SystemdService interface {
	// Name sets the name of the service.
	Name(name string) SystemdService
	// Root sets the root path of the service.
	Root(dir string) SystemdService
	// Command sets the command of the service.
	Command(command string) SystemdService
	// Template sets the template for the service.
	// template string can contain {name}, {root} and {command} placeholders.
	Template(engine TemplateEngine) SystemdService
	// Exists checks if the service exists.
	Exists() bool
	// Enabled checks if the service exists and enabled on startup.
	Enabled() bool
	// Install installs the service.
	// override parameter indicating whether to override existing configurations.
	// returns false if service exists and not override.
	Install(override bool) (bool, error)
	// Uninstall uninstalls the service.
	Uninstall() error
}

type systemdDriver struct {
	name     string
	root     string
	command  string
	template TemplateEngine
}

func (driver systemdDriver) path() string {
	return "/etc/systemd/system/" + driver.name + ".service"
}

func (driver *systemdDriver) Name(name string) SystemdService {
	driver.name = name
	return driver
}

func (driver *systemdDriver) Root(dir string) SystemdService {
	driver.root = dir
	return driver
}

func (driver *systemdDriver) Command(command string) SystemdService {
	driver.command = command
	return driver
}

func (driver *systemdDriver) Template(engine TemplateEngine) SystemdService {
	driver.template = engine
	return driver
}

func (driver *systemdDriver) Exists() bool {
	_, err := exec.Command("sudo", "systemctl", "status", driver.name).Output()
	return err == nil
}

func (driver *systemdDriver) Enabled() bool {
	output, _ := exec.Command("sudo", "systemctl", "is-enabled", driver.name).Output()
	return strings.HasPrefix(string(output), "enabled")
}

func (driver *systemdDriver) Install(override bool) (bool, error) {
	if exists := driver.Exists(); exists && !override {
		return false, nil
	}

	content := driver.template.
		AddParameter("name", driver.name).
		AddParameter("root", driver.root).
		AddParameter("command", driver.command).
		Compile()

	err := os.WriteFile(driver.path(), []byte(content), 0644)
	if err != nil {
		return false, err
	}

	err = exec.Command("sudo", "systemctl", "daemon-reload").Run()
	if err != nil {
		return false, err
	}

	err = exec.Command("sudo", "systemctl", "enable", driver.name).Run()
	if err != nil {
		return false, err
	}

	err = exec.Command("sudo", "systemctl", "start", driver.name).Run()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (driver *systemdDriver) Uninstall() error {
	if driver.Exists() {
		err := exec.Command("systemctl", "stop", driver.name).Run()
		if err != nil {
			return err
		}

		err = exec.Command("systemctl", "disable", driver.name).Run()
		if err != nil {
			return err
		}
	}

	return os.Remove(driver.path())
}
