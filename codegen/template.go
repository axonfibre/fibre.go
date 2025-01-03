package codegen

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/axonfibre/fibre.go/ierrors"
	"github.com/axonfibre/fibre.go/lo"
)

// Template is a wrapper around the text/template package that provides a generic way for generating files according
// to the "go generate" pattern.
//
// In addition to the standard template delimiters (https://pkg.go.dev/text/template), it supports /*{{ ... }}*/, which
// allows to "hide" template code in go comments. Data is provided to the template as function pipelines.
type Template struct {
	// header contains the "fixed" header of the file above the "go generate" statement (not processed by the template).
	header string

	// content contains "dynamic" the content of the file below the "go generate" statement.
	content string

	// mappings is a set of tokens that are being mapped to pipelines in the template.
	mappings template.FuncMap
}

// NewTemplate creates a new Template with the given pipeline mappings.
func NewTemplate(mappings template.FuncMap) *Template {
	return &Template{
		mappings: mappings,
	}
}

// Parse parses the given file and extracts the header and content by splitting the file at the "go:generate" directive.
// It automatically removes existing "//go:build ignore" directives from the header.
func (t *Template) Parse(fileName string) error {
	readFile, err := os.ReadFile(fileName)
	if err != nil {
		return ierrors.Wrapf(err, "could not read file %s", fileName)
	}

	splitTemplate := strings.Split(string(readFile), "//go:generate")
	if len(splitTemplate) != 2 {
		return ierrors.Errorf("could not find go:generate directive in %s", fileName)
	}

	t.header = strings.TrimSpace(strings.ReplaceAll(splitTemplate[0], "//go:build ignore", ""))
	t.content = strings.TrimSpace(splitTemplate[1][strings.Index(splitTemplate[1], "\n"):])

	return nil
}

// Generate generates the file with the given fileName (it can receive an optional generator function that overrides the
// way the content is generated).
func (t *Template) Generate(fileName string, optGenerator ...func() (string, error)) error {
	generatedContent, err := lo.First(optGenerator, t.GenerateContent)()
	if err != nil {
		return ierrors.Wrap(err, "could not generate content")
	}

	//nolint:gosec // false positive, only used for code generation
	return os.WriteFile(fileName, []byte(strings.Join([]string{
		generatedFileHeader + t.header,
		generatedContent + "\n",
	}, "\n\n")), 0644)
}

// GenerateContent generates the dynamic content of the file by processing the template.
func (t *Template) GenerateContent() (string, error) {
	// replace /*{{ and }}*/ with {{ and }} to "unpack" statements that are embedded as comments
	content := regexp.MustCompile(`/\*{{`).ReplaceAll([]byte(t.content), []byte("{{"))
	content = regexp.MustCompile(`}}\*/`).ReplaceAll(content, []byte("}}"))

	tmpl, err := template.New("template").Funcs(t.mappings).Parse(string(content))
	if err != nil {
		return "", ierrors.Wrap(err, "could not parse template")
	}

	buffer := new(bytes.Buffer)
	if err := tmpl.Execute(buffer, nil); err != nil {
		return "", ierrors.Wrap(err, "could not execute template")
	}

	return buffer.String(), nil
}

// generatedFileHeader is the header that is being added to the top of the generated file.
const generatedFileHeader = "// Code generated by go generate; DO NOT EDIT.\n"
