package main

import (
	"askme/pkg/models"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)



func (app *application) logPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", nil)
}
func (app *application) signUpPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", nil)
}

func (app *application) adminPage(w http.ResponseWriter, r *http.Request) {

	td :=&templateData{}
	templ := "admin.page.tmpl"
	if app.isAdmin(r) {
		templ= "admin.page.tmpl"
		td.Questions,_=app.question.GetAllQuestions()

	}else {
		templ= "home.page.tmpl"
	}
	app.render(w, r, templ, td)
}

func (app application) answer(w http.ResponseWriter,r *http.Request){
	if !app.isAuthenticated(r) {
		qs,_:=app.question.GetAllQuestions()
		app.render(w,r,"home.page.tmpl",&templateData{
			Questions: qs,
			IsError: true,
			Not: "please login or signUP",
		})
	}
	err:=r.ParseForm()

	if err!=nil {
		app.serverError(w,err)
	}
	q_id,_:=strconv.Atoi(r.PostForm.Get("q_id"))
	ans:= r.PostForm.Get("answer")
	fmt.Println(q_id,ans)
	person,err:=app.getPFS(r)
	app.question.AnswerQuestion(person.Id,q_id,ans)


	http.Redirect(w,r,"/",http.StatusSeeOther)




}
func (app application) askQuestion(w http.ResponseWriter,r *http.Request){
	if !app.isAuthenticated(r) {
		qs,_:=app.question.GetAllQuestions()
		app.render(w,r,"home.page.tmpl",&templateData{
			Questions: qs,
			IsError: true,
			Not: "please login or signUP",
		})
	}
	person,err:=app.getPFS(r)
	if err!=nil {

		fmt.Println(person)
	}
	err=r.ParseForm()
	q:=r.PostForm.Get("question")
	err=app.question.AskQuestion(person.Id,q)
	if err!=nil {
		app.serverError(w,err)
	}
	app.home(w,r)



}
func (app application) Question(w http.ResponseWriter,r *http.Request){
	idd := r.URL.Query().Get("id")
	id,_:=strconv.Atoi(idd)
	q,err:=app.question.GetSingleQuestion(id)

	if err!=nil {
		app.serverError(w,err)
	}
	fmt.Println(q.Quest)
	app.render(w,r,"show.page.tmpl",&templateData{
		Question: q,
	})

}

func (app application) deleteQuestion(w http.ResponseWriter,r *http.Request){
	err:=r.ParseForm()
	if err!=nil {
		app.serverError(w,err)
	}

	id,_:=strconv.Atoi(r.PostForm.Get("q_id"))
	app.question.DeleteQuestion(id)
	app.adminPage(w,r)

}

func (app application) deleteAnswer(w http.ResponseWriter,r *http.Request){
	err:=r.ParseForm()
	if err!=nil {
		app.serverError(w,err)
	}

	id,_:=strconv.Atoi(r.PostForm.Get("a_id"))
	app.question.DeleteAnswer(id)
	app.adminPage(w,r)
}
func (app application) getPFS(r *http.Request) (*models.User,error) {
	sess,_:=Store.Get(r,"auth")
	val := sess.Values["id"]
	pers,err:=app.persons.GetPersonBylogin(val.(string))
	if err!=nil {
		return nil,err
	}
	return pers,nil
}















func (app *application)login(w http.ResponseWriter, r *http.Request)   {

	td :=&templateData{}
	err:=r.ParseForm()
	if err!=nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err=nil
	login := r.PostForm.Get("login")
	pass := r.PostForm.Get("password")
	addmin := r.Form["asAdmin"]

	pers:=&models.User{}
	if addmin!=nil{
		pers,err = app.persons.GetAdmin(login,pass)
		fmt.Println(err, "admin err")
		if err==nil {
			fmt.Println(err)
			td.Admin=true
			isAdmin,_:= Store.Get(r,"admin")
			isAdmin.Values["adm"]=pers.Name
			isAdmin.Save(r,w)
		}else {
			td.Admin=false
		}
	}else {
		fmt.Println("ne admin")
		pers,err = app.persons.GetPerson(login,pass)
	}

	var templ string
	if err!=nil{
		td.IsError =true
		td.Not ="wrong password or login"
		app.serverError(w,err)
		templ="login.page.tmpl"
		app.render(w,r,templ,td)
	}else {

		sess,_:=Store.Get(r,"auth")
		sess.Values["id"]=pers.Login
		sess.Save(r,w)
		http.Redirect(w,r,"/",http.StatusSeeOther)

	}




}
func (app *application) logOut(w http.ResponseWriter, r *http.Request)  {
	Sess:=sessions.NewSession(Store,"auth")
	Sess.Options.MaxAge=-1
	Sess.Save(r,w)
	isAdmin:= sessions.NewSession(Store,"admin")
	isAdmin.Options.MaxAge=-1
	isAdmin.Save(r,w)
	http.Redirect(w,r,"/",http.StatusSeeOther)
}
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	qs ,err:=app.question.GetAllQuestions()
	
	if err!=nil {

		app.serverError(w,err)
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Questions: qs,
	})



}

func (app *application)signUP(w http.ResponseWriter, r *http.Request) {

	td := &templateData{}
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)

	}
	err = nil
	login := r.PostForm.Get("login")
	pass := r.PostForm.Get("password")
	name := r.PostForm.Get("name")
	phone := r.PostForm.Get("phone")
	pers := &models.User{Login: login, Passrword: pass, Name: name, Phone:phone}
	p,_:=app.persons.GetPersonBylogin(pers.Login)
	var templ string
	if p==nil{
		err = app.persons.SignUp(pers)
		templ="home.page.tmpl"
		td.User = pers
		td.IsError =false
		sess,_:=Store.Get(r,"auth")
		sess.Values["id"]=pers.Login
		sess.Save(r,w)
		http.Redirect(w,r,"/",http.StatusSeeOther)


	}else{
		td.IsError =true
		td.Not ="this user is alredy exist"
		app.serverError(w,err)
		templ="signup.page.tmpl"
	}
	app.render(w,r,templ,td)



}
