package main

import (
	"bytes"
	"fmt"
	"net/http"
)


func (app *application) serverError(w http.ResponseWriter, err error) {
	/*trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)*/
}
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}


func (app *application) isAuthenticated(r *http.Request) bool {
	sess,_ := Store.Get(r,"auth")

	if sess.IsNew {
		return false
	}
	return true
}
func (app *application) isAdmin(r *http.Request) bool {
	sess,_ := Store.Get(r,"admin")

	if sess.IsNew {
		return false
	}
	return true
}
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {

	if td == nil {
		td = &templateData{}
	}
	td.Authenticated = app.isAuthenticated(r)
	td.Admin=app.isAdmin(r)


	if app.isAuthenticated(r) {

		sess,_:=Store.Get(r,"auth")
		val := sess.Values["id"]

		id:= val.(string)

		td.User,_=app.persons.GetPersonBylogin(id)

	}


	return td
}
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		fmt.Println(err)

		app.serverError(w, err)
		return
	}



	buf.WriteTo(w)


}