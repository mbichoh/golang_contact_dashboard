package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mbichoh/contactDash/pkg/models/mysql"
)

type application struct {
	errorLog        *log.Logger
	infoLog         *log.Logger
	session         *sessions.Session
	contacts        *mysql.ContactModel
	templateCache   map[string]*template.Template
	users           *mysql.UserModel
	groups          *mysql.GroupsModel
	groupedcontacts *mysql.GroupedContactsModel
}

type contextKey string

var contextKeyUser = contextKey("user")

func main() {

	addr := flag.String("addr", ":7076", "	HTTP network address")
	dsn := flag.String("dsn", "golang:goconnect@/demogo?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:        errorLog,
		infoLog:         infoLog,
		session:         session,
		contacts:        &mysql.ContactModel{DB: db},
		users:           &mysql.UserModel{DB: db},
		groups:          &mysql.GroupsModel{DB: db},
		groupedcontacts: &mysql.GroupedContactsModel{DB: db},
		templateCache:   templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Server starting on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
