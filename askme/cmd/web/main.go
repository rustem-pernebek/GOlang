package main

import (
	"askme/pkg/models/postgres"
	"context"

	"flag"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"

)

var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type application struct {
	/*errorLog *log.Logger*/
	infoLog *log.Logger

	persons *postgres.PersonModel
	session *sessions.Session
	question *postgres.QuestionModel

	templateCache map[string]*template.Template
}



func main() {
	Store.Options.Domain="localhost"
	addr := flag.String("addr", ":4001", "HTTP network address")
	flag.Parse()


	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=12345 host=localhost port=4000 dbname=postgres sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}

	defer pool.Close()
	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}




	app := &application{
		/*errorLog:errorLog,*/
		infoLog:infoLog,

		persons :&postgres.PersonModel{ pool},

		question: &postgres.QuestionModel{pool},
		templateCache: templateCache}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}