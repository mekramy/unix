package unix

import "strings"

func NewEngine() TemplateEngine {
	engine := new(engineDriver)
	engine.params = make([]string, 0)
	return engine
}

type TemplateEngine interface {
	SetTemplate(template string) TemplateEngine
	AddParameter(name, value string) TemplateEngine
	Compile() string
}

type engineDriver struct {
	template string
	params   []string
}

func (e *engineDriver) SetTemplate(template string) TemplateEngine {
	e.template = template
	return e
}

func (e *engineDriver) AddParameter(name, value string) TemplateEngine {
	e.params = append(e.params, "{"+name+"}", value)
	return e
}

func (e *engineDriver) Compile() string {
	return strings.NewReplacer(e.params...).Replace(e.template)
}
