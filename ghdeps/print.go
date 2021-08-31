package ghdeps

import "html/template"

type (
	PrintOption struct {
		Template *template.Template
	}
)

var (
	DefaultTemplate = template.Must(template.New("default").Parse(
		`Dependents of {{.Source.User}}/{{.Source.Repo}}
----------------------------
TOTAL:	{{len .Dependents}}
PAGES:	{{len .Pages}}
----------------------------
{{range .Dependents}}{{.User}}/{{.Repo}}
{{end}}`,
	))
)

func (opt *PrintOption) ensure() *PrintOption {
	if opt == nil {
		return &PrintOption{
			Template: DefaultTemplate,
		}
	}
	if opt.Template == nil {
		opt.Template = DefaultTemplate
	}
	return opt
}
