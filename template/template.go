// Package template simplifies parsing a set of templates with common layouts.
package template

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

// MissingError is returned when trying to render a template that does not
// exist.
type MissingError string

func (e MissingError) Error() string {
	return fmt.Sprintf("no template exists with name '%s'", string(e))
}

type Templates map[string]*template.Template

// Parse each html/template in templateDir with those in templateDir/layout,
// providing the custom funcs.
func Parse(templateDir string, funcs template.FuncMap) (Templates, error) {
	layouts, err := template.New("").Funcs(funcs).ParseGlob(filepath.Join(templateDir, "layout/*.*"))
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(templateDir, "*.*"))
	if err != nil {
		return nil, err
	}

	tmpls := map[string]*template.Template{}
	for _, file := range files {
		clone, err := layouts.Clone()
		if err != nil {
			return nil, err
		}

		tmpl, err := clone.ParseFiles(file)
		if err != nil {
			return nil, err
		}

		tmpls[filepath.Base(file)] = tmpl
	}

	return tmpls, nil
}

type Template func(io.Writer, interface{}) error

// Get name template from the Templates collection. It assumes that the defined
// template to be rendered is called "page".
func (t Templates) Get(name string) Template {
	tmpl := t[name]

	if tmpl == nil {
		return func(wr io.Writer, data interface{}) error {
			return MissingError(name)
		}
	}

	return func(wr io.Writer, data interface{}) error {
		return tmpl.ExecuteTemplate(wr, "page", data)
	}
}
