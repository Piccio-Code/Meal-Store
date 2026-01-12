package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	config   config
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "The port of the backend.")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.LstdFlags|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		config:   cfg,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.infoLog.Printf("starting %s server on http://localhost:%d", cfg.env, cfg.port)
	app.errorLog.Fatal(srv.ListenAndServe())
}
