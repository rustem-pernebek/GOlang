package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)


func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)


	mux := pat.New()
	mux.Get("/askme/question/",http.HandlerFunc(app.Question))
	mux.Get("/askme/questions",http.HandlerFunc(app.Question))
	mux.Post("/askme/answer",http.HandlerFunc(app.answer))
	mux.Post("/askme/ask",http.HandlerFunc(app.askQuestion))
	mux.Post("/askme/deleteQuestion",http.HandlerFunc(app.deleteQuestion))
	mux.Post("/askme/deleteAnswer",http.HandlerFunc(app.deleteAnswer))
	mux.Get("/",http.HandlerFunc(app.home))
	mux.Get("/askme/login",http.HandlerFunc(app.logPage))
	mux.Post("/askme/login",http.HandlerFunc(app.login))
	mux.Get("/logout", http.HandlerFunc(app.logOut))
	mux.Get("/admin",http.HandlerFunc(app.adminPage))


	mux.Get("/askme/signup",http.HandlerFunc(app.signUpPage))
	mux.Post("/askme/signup",http.HandlerFunc(app.signUP))



	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}