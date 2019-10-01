package main

import (
	"net/http"

	"github.com/bmizerany/pat"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, app.authenticate)

	mux := pat.New()

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signup))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.login))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logout))

	mux.Get("/", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.ContactHome))
	mux.Get("/contact/group", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.DispGroupedContacts))
	mux.Post("/contact/group", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.CreateGroupedContacts))

	mux.Get("/contact/group/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.GroupedContacts))
	mux.Get("/contact/del/usergroup/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.DelGroupContact))
	mux.Get("/contact/del/group/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.DelGroup))
	mux.Post("/contact/del/group/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.SendMessageToGroup))
	mux.Get("/contact/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.CreateContactForm))
	mux.Post("/contact/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.CreateContact))
	mux.Get("/contact/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.ShowContact))
	mux.Post("/contact/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.SendMessageToContact))
	mux.Get("/contact/update/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.FetchUpdateContact))
	mux.Post("/contact/update/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.UpdateContact))
	mux.Get("/contact/del/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.DelContact))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
