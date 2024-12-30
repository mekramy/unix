package unix

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func NewNginxReverseProxy(name, port string) ServerBlock {
	server := new(serverBlock)
	server.name = name
	server.port = port
	server.template = NewEngine()
	server.template.SetTemplate(`
server {
        listen 80;
        listen [::]:80;
        server_name {domains};

        location / {
            client_max_body_size 1M;
            proxy_pass http://localhost:{port};
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header Referer $http_referer;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_set_header X-Forwarded-For $remote_addr;
            proxy_set_header X-Forwarded-Referer $http_referer;
            proxy_cache_bypass $http_upgrade;
        }
}
	`)
	return server
}

type ServerBlock interface {
	// Name sets the name of the site.
	Name(name string) ServerBlock
	// Port sets the port for the site.
	Port(port string) ServerBlock
	// Domains sets the domains for the site.
	Domains(domains ...string) ServerBlock
	// Template sets the template for the site.
	// template string can contain {domains} and {port} placeholders.
	Template(engine TemplateEngine) ServerBlock
	// Disable disables the site manually.
	Disable() error
	// Enable enables the site manually.
	Enable() error
	// Exists checks if the site exists.
	Exists() (bool, error)
	// Enabled checks if the site exists and enabled.
	Enabled() (bool, error)
	// Install installs the site.
	// override parameter indicating whether to override existing configurations.
	// returns false if site exists and not override.
	Install(override bool) (bool, error)
	// Uninstall uninstalls the site.
	Uninstall() error
}

type serverBlock struct {
	name     string
	domains  []string
	port     string
	template TemplateEngine
}

func (server serverBlock) path() string {
	return "/etc/nginx/sites-available/" + server.name
}

func (server serverBlock) link() string {
	return "/etc/nginx/sites-enabled/" + server.name
}

func (server *serverBlock) Name(name string) ServerBlock {
	server.name = name
	return server
}

func (server *serverBlock) Port(port string) ServerBlock {
	server.port = port
	return server
}

func (server *serverBlock) Domains(domains ...string) ServerBlock {
	server.domains = append(server.domains, domains...)
	return server
}

func (server *serverBlock) Template(engine TemplateEngine) ServerBlock {
	server.template = engine
	return server
}

func (server *serverBlock) Disable() error {
	// delete link
	if exists, err := FileExists(server.link()); err != nil {
		return err
	} else if !exists {
		return nil
	} else if err := os.Remove(server.link()); err != nil {
		return err
	}

	// Reload nginx to apply the changes
	return exec.Command("systemctl", "restart", "nginx").Run()
}

func (server *serverBlock) Enable() error {
	// delete link
	if exists, err := FileExists(server.path()); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("%s file not exists", server.path())
	} else if err := os.Symlink(server.path(), server.link()); err != nil {
		return err
	}

	// Reload nginx to apply the changes
	return exec.Command("systemctl", "restart", "nginx").Run()
}

func (server *serverBlock) Exists() (bool, error) {
	return FileExists(server.path())
}

func (server *serverBlock) Enabled() (bool, error) {
	if available, err := FileExists(server.path()); err != nil {
		return false, err
	} else if enabled, err := FileExists(server.link()); err != nil {
		return false, err
	} else {
		return available && enabled, nil
	}
}

func (server *serverBlock) Install(override bool) (bool, error) {
	content := server.template.
		AddParameter("port", server.port).
		AddParameter("domains", strings.Join(server.domains, " ")).
		Compile()

	if exists, err := FileExists(server.path()); err != nil {
		return false, err
	} else if exists && !override {
		return false, nil
	}

	if err := os.WriteFile(server.path(), []byte(content), 0644); err != nil {
		return false, err
	}

	if err := os.Symlink(server.path(), server.link()); err != nil {
		return false, err
	}

	if err := exec.Command("systemctl", "restart", "nginx").Run(); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (server *serverBlock) Uninstall() error {
	// Remove the enabled site link
	if err := os.Remove(server.link()); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Remove the available site file
	if err := os.Remove(server.path()); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Reload nginx to apply the changes
	return exec.Command("systemctl", "restart", "nginx").Run()
}
