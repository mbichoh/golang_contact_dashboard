package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mbichoh/contactDash/pkg/models"
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

	// CHECK : You should have one config for all this, port, passwords, etc

	config := new(models.Config)

	flag.StringVar(&config.Addr, "addr", models.AddPort, "HTTP network address")
	flag.StringVar(&config.DSN, "dsn", models.Dsn, "MySQL data source name")
	flag.StringVar(&config.Secret, "secret", models.SecretKey, "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*&config.DSN)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*&config.Secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

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

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *&config.Addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Server starting on %s", *&config.Addr)

	// CHECK : I dont see any tls, install local tls certificate and ListenAndServeTLS('path'/to/tls)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
