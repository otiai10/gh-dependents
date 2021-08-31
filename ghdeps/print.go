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
	JSONTemplate = template.Must(template.New("json").Parse(
		`{
    "source": {
        "user": "{{.Source.User}}",
        "repo": "{{.Source.Repo}}",
        "url": "{{.Source.URL .ServiceURL}}"
    },
    "dependents": [{{range $i, $d := .Dependents}}{{if $i}},{{end}}
        {
            "user": "{{$d.User}}",
            "repo": "{{$d.Repo}}",
            "url": "{{$d.URL $.ServiceURL}}"
        }{{end}}
    ],
    "pages": [{{range $i, $p := .Pages}}{{if $i}},{{end}}
        "{{$p}}"{{end}}
    ]
}`,
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
