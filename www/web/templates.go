package main

import (
	"html/template"
	"path/filepath"

	"github.com/mbichoh/contactDash/pkg/forms"
	"github.com/mbichoh/contactDash/pkg/models"
)

type templateData struct {
	AuthenticatedUser *models.User
	CSRFToken         string
	Flash             string
	Form              *forms.Form
	User              *models.User
	Contact           *models.Contact
	Contacts          []*models.Contact
	Group             *models.Groups
	Groups            []*models.Groups
	GroupedContacts   *models.GroupedContacts
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
