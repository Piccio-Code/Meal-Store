package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port       int
	env        string
	poolConfig poolConfig
}

type poolConfig struct {
	MaxConns          int32         `json:"maxConns,omitempty"`
	MinConns          int32         `json:"minConns,omitempty"`
	MaxConnLifetime   time.Duration `json:"maxConnLifetime,omitempty"`
	MaxConnIdleTime   time.Duration `json:"maxConnIdleTime,omitempty"`
	HealthCheckPeriod time.Duration `json:"healthCheckPeriod,omitempty"`
}

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	pool     *pgxpool.Pool
	config   config
}

func main() {
	poolConfig := poolConfig{
		MaxConns:          10,
		MinConns:          0,
		MaxConnLifetime:   time.Hour * 1,
		MaxConnIdleTime:   time.Minute * 30,
		HealthCheckPeriod: time.Minute,
	}

	cfg := config{
		poolConfig: poolConfig,
	}

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

	pool.Config().MaxConns = cfg.poolConfig.MaxConns
	pool.Config().MinConns = cfg.poolConfig.MinConns
	pool.Config().MaxConnLifetime = cfg.poolConfig.MaxConnLifetime
	pool.Config().MaxConnIdleTime = cfg.poolConfig.MaxConnIdleTime
	pool.Config().HealthCheckPeriod = cfg.poolConfig.HealthCheckPeriod

	defer pool.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		pool:     pool,
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

func (p poolConfig) MarshalJSON() ([]byte, error) {
	type poolConfigOut struct {
		MaxConns          int32  `json:"maxConns,omitempty"`
		MinConns          int32  `json:"minConns,omitempty"`
		MaxConnLifetime   string `json:"maxConnLifetime,omitempty"`
		MaxConnIdleTime   string `json:"maxConnIdleTime,omitempty"`
		HealthCheckPeriod string `json:"healthCheckPeriod,omitempty"`
	}

	js := poolConfigOut{
		MaxConns:          p.MaxConns,
		MinConns:          p.MinConns,
		MaxConnLifetime:   p.MaxConnLifetime.String(),
		MaxConnIdleTime:   p.MaxConnIdleTime.String(),
		HealthCheckPeriod: p.HealthCheckPeriod.String(),
	}

	return json.MarshalIndent(js, "", "\t")
}
