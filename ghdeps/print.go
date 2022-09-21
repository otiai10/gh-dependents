package ghdeps

import "html/template"

type (
	PrintOption struct {
		Template *template.Template
		Sort     SortProvider
	}

	SortProvider func(deps Dependents) func(int, int) bool
)

var (
	SortByStar SortProvider = func(deps Dependents) func(int, int) bool {
		return func(i, j int) bool {
			return deps[i].Stars > deps[j].Stars
		}
	}
	SortByFork SortProvider = func(deps Dependents) func(int, int) bool {
		return func(i, j int) bool {
			return deps[i].Forks > deps[j].Forks
		}
	}
)

var (
	PrettyTemplate = template.Must(template.New("pretty").Parse(
		`--------------------------------------
Dependents of {{.Source.User}}/{{.Source.Repo}}
TOTAL:	{{len .Dependents}}
PAGES:	{{len .Pages}}
--------------------------------------
{{range .Dependents}}⭐️ {{.Stars}}	🌵 {{.Forks}}	{{.User}}/{{.Repo}}
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
    "page":  {{if eq .PageCount 0}}null{{else}}{{.PageCount}}{{end}},
    "after": {{if eq (len .After) 0}}null{{else}}"{{.After}}"{{end}}
  },
  "dependents": [{{range $i, $d := .Dependents}}{{if $i}},{{end}}
    {
      "user": "{{$d.User}}",
      "repo": "{{$d.Repo}}",
      "url": "{{$d.URL $.ServiceURL}}",
      "stars": {{$d.Stars}},
      "forks": {{$d.Forks}}
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
			Template: PrettyTemplate,
		}
	}
	if opt.Template == nil {
		opt.Template = PrettyTemplate
	}

	return opt
}
