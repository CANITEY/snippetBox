package main

import (
	"caniteySnippetBox/internal/models"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"html/template"

	_ "github.com/go-sql-driver/mysql"
)


type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *models.SnippetModel
	templateCache map[string]*template.Template
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

	app := &application{
		errorLog,
		infoLog,
		&models.SnippetModel{DB:db},
		templateCache,
	}



	srv := &http.Server{
		Addr: *address,
		ErrorLog: errorLog,
		Handler: app.routes(),
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
