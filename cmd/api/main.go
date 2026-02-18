package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Piccio-Code/MealStore/internal/data"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	pool     *pgxpool.Pool
	models   data.Models
	config   config
}

func main() {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	flag.IntVar(&cfg.port, "port", 8080, "The port of the backend.")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.LstdFlags|log.Lshortfile)

	pool, err := NewDBPool()
	if err != nil {
		log.Println("Error connecting to the DB")
		log.Fatal(err)
	}

	defer pool.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		pool:     pool,
		models:   data.NewModels(pool),
		config:   cfg,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.infoLog.Printf("starting %s server on http://localhost:%d", cfg.env, cfg.port)
	app.errorLog.Fatal(srv.ListenAndServe())
}

func NewDBPool() (pool *pgxpool.Pool, err error) {
	pool, err = pgxpool.New(context.Background(), os.Getenv("DB_DSN"))
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
