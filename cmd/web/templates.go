package main

import (
	"caniteySnippetBox/internal/models"
	"caniteySnippetBox/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
	Form any
	Flash string
	IsAuthenticated bool
	CSRFToken	string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jun 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}


	for _, page := range pages {
		name := filepath.Base(page)


		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, "html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFS(ui.Files, "html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}


		ts, err = ts.ParseFS(ui.Files, page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
