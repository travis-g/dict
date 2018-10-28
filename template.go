package main

import (
	"strings"
	"text/template"
)

// Helper function to capitalize the first character of a string
func capitalizeFirst(s string) string {
	if len(s) > 1 {
		return strings.ToUpper(string(s[0])) + s[1:]
	} else if len(s) == 1 {
		return strings.ToUpper(string(s[0]))
	} else {
		return ""
	}
}

// Helper function to split a string on an old delimeter and re-join the string
// slice with a new one.
func rejoin(input, old, new string) string {
	input = strings.TrimRight(input, old)
	tmp := strings.Split(input, old)
	return strings.Join(tmp, new)
}

// Send functions to the template rendering engine
var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"join": func(arr []string, sep string) string {
		return strings.Join(arr, sep)
	},
	"rejoin": rejoin,
}

// Templates is the global templates used to render API results.
var Templates = map[string]*template.Template{
	"definition": template.Must(template.New("definition").Funcs(funcMap).Parse(`
{{- .RenderTitle }}

{{ range .LexicalEntries -}}
{{ .RenderLexicalCategory }}:

{{ if len .Pronunciations }}{{ range .Pronunciations }}- {{ .String }}
{{ end }}
{{ end -}}
{{ range .Entries -}}
{{ range $i, $sense := .Senses -}}
{{inc $i}}. {{if .Tags }}{{ .RenderTags }} {{end}}{{ .RenderDefinitions }}
{{ range .Examples }}   - {{ .Render }}
{{ end }}
{{- if .Subsenses }}
{{ range $j, $subsense := .Subsenses }}   {{inc $j}}. {{if .Tags }}{{ .RenderTags }} {{end}}{{ .RenderDefinitions }}
{{ range .Examples }}      - {{ .Render }}
{{ end }}
{{ end -}}{{ else }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}
`)),
	"definition-simple": template.Must(template.New("definition-simple").Funcs(funcMap).Parse(`
{{- range .LexicalEntries -}}
{{ range .Entries -}}
{{ range .Senses -}}
{{ .RenderDefinitions }}
{{ end -}}
{{ end -}}
{{ end -}}
`)),
	"synonyms": template.Must(template.New("synonym").Funcs(funcMap).Parse(`
{{- .RenderTitle }}

{{ range .LexicalEntries -}}
{{ .RenderLexicalCategory }}:
{{ range .Entries }}
{{ range $i, $sense := .Senses }}{{inc $i}}. {{ rejoin .RenderExamples "\n" ", " }}
{{- if .HasSynonyms }}
   - {{ .RenderTags "informal" }}{{ .RenderSynonyms }}
{{- end -}}
{{- if .Subsenses }}
{{ range .Subsenses -}}
{{- if .HasSynonyms }}   - {{if .Tags "informal"}}{{ .RenderTags "informal" }} {{end}}{{ .RenderSynonyms }}
{{ end -}}
{{ end -}}{{ else }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}
`)),
	"synonyms-simple": template.Must(template.New("synonyms-simple").Funcs(funcMap).Parse(`
{{- range .LexicalEntries -}}
{{ range .Entries -}}
{{ range .Senses -}}
{{- if .HasSynonyms }}{{ .RenderSynonyms }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}
`)),
	"antonyms": template.Must(template.New("antonym").Funcs(funcMap).Parse(`
{{- .RenderTitle }}

{{ range .LexicalEntries -}}
{{ .RenderLexicalCategory }}:
{{ range .Entries }}
{{ range $i, $sense := .Senses }}{{inc $i}}. {{ rejoin .RenderExamples "\n" ", " }}
{{- if .HasAntonyms }}
   - {{ .RenderTags "informal" }}{{ .RenderAntonyms }}
{{- end -}}
{{- if .Subsenses }}
{{ range .Subsenses -}}
{{- if .HasAntonyms }}   - {{if .Tags "informal"}}{{ .RenderTags "informal" }} {{end}}{{ .RenderAntonyms }}
{{ end -}}
{{ end -}}{{ else }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}
`)),
	"antonyms-simple": template.Must(template.New("antonyms-simple").Funcs(funcMap).Parse(`
{{- range .LexicalEntries -}}
{{ range .Entries -}}
{{ range .Senses -}}
{{- if .HasAntonyms }}{{ .RenderAntonyms }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}
`)),
	"webpage": template.Must(template.New("webpage").Funcs(funcMap).Parse(`<!DOCTYPE html>
<html lang="en">
{{if .Title}}<title>{{.Title}}</title>{{end}}
<script>window.mdme = {style: 'none'}</script>
<script src="https://unpkg.com/mdme"></script>
<style>html{font-family:sans-serif}</style>
<textarea>
{{ .Content -}}
</textarea>
</html>
`)),
}
