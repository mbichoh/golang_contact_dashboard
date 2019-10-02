package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/mbichoh/contactDash/pkg/models"
)

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.session.PopString(r, "flash")
	td.AuthenticatedUser = app.authenticatedUser(r)
	return td
}

//error 500 "Internal Server Error"
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n", err.Error())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//error 400 "Bad Request"
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

//error 404 "Not Found"
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s doesn't exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}
func (app *application) authenticatedUser(r *http.Request) *models.User {

	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}

	return user

}
