package ghdeps

import "html/template"

type (
	PrintOption struct {
		Template   *template.Template
		SortByStar bool
	}
)

var (
	DefaultTemplate = template.Must(template.New("default").Parse(
		`--------------------------------------
Dependents of {{.Source.User}}/{{.Source.Repo}}
TOTAL:	{{len .Dependents}}
PAGES:	{{len .Pages}}
--------------------------------------
{{range .Dependents}}⭐️{{.Stars}}	{{.User}}/{{.Repo}}
{{end}}`,
	))

	JSONTemplate = template.Must(template.New("json").Parse(
		`{
    "source": {
        "user": "{{.Source.User}}",
        "repo": "{{.Source.Repo}}",
        "url": "{{.Source.URL .ServiceURL}}"
    },
    "query": {
        "page":  {{if eq .PageNum 0}}null{{else}}{{.PageNum}}{{end}},
        "after": {{if eq (len .After) 0}}null{{else}}"{{.After}}"{{end}}
    },
    "dependents": [{{range $i, $d := .Dependents}}{{if $i}},{{end}}
        {
            "user": "{{$d.User}}",
            "repo": "{{$d.Repo}}",
            "url": "{{$d.URL $.ServiceURL}}",
            "stars": {{$d.Stars}}
        }{{end}}
    ],
    "pages": [{{range $i, $p := .Pages}}{{if $i}},{{end}}
        {
            "url":  "{{$p.URL}}"{{if ne (len $p.Next) 0}},
            "next": "{{$p.Next}}"{{end}}
        }{{end}}
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
