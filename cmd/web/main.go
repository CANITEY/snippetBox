package main

import (
	"caniteySnippetBox/internal/models"
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)


type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	address := flag.String("addr", ":8888", "HTTP address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	// database connection
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// template caching utility
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// create a session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog,
		infoLog,
		&models.SnippetModel{DB:db},
		templateCache,
		form.NewDecoder(),
		sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}



	srv := &http.Server{
		Addr: *address,
		ErrorLog: errorLog,
		Handler: app.routes(),
		TLSConfig: tlsConfig,
	}
	infoLog.Printf("Serving at %s", *address)

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
