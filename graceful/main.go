package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Variables injected by -X flag
var AppVersion = "unknown"
var GitVersion = "unknown"
var LastCommitTime = "unknown"
var LastCommitHash = "unknown"
var LastCommitUser = "unknown"
var BuildTime = "unknown"

// Timeout is the duration to allow outstanding requests to survive before forcefully terminating them.
const Timeout = 20

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("can't load configuration: %s", err)
	}

	db, err := sql.Open("mysql", cfg.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	app, err := NewAppRest(db)
	if err != nil {
		log.Fatal(err)
	}

	users := NewUserRest(db)

	router := mux.NewRouter()
	router.HandleFunc("/api/users", users.Users).Methods("GET")
	router.HandleFunc("/api/users", users.AddUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", users.GetUser).Methods("GET")
	router.HandleFunc("/ready", app.Ready).Methods("GET")
	router.HandleFunc("/health", app.Health).Methods("GET")
	router.HandleFunc("/api/error", app.Err).Methods("POST")
	router.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))

	log.Println("starting http server with graceful shutdown mode!")

	// create and start http server in new goroutine
	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.HTTPPort), Handler: router}
	go func() {
		// we can't use log.Fatal here!
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("http server stoped: %s\n", err)
		}
	}()

	// subscribe to SIGTERM signals
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// blocks the execution until os.Interrupt or syscall.SIGTERM signal appears
	<-quit
	log.Println("shutting down server. waiting to drain the ongoing requests...")
	app.Unhealthy()

	// add extra time to prevent new requests be routed to our service
	time.Sleep(5 * time.Second)

	// shut down gracefully, but wait no longer than the Timeout value.
	ctx, cancelF := context.WithTimeout(context.Background(), Timeout*time.Second)
	defer cancelF()

	// shutdown the http server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error while shutdown http server: %v\n", err)
	}

	log.Println("server gracefully stopped")
}

type config struct {
	HTTPPort   int    `split_words:"true" default:"8080"`
	DBHost     string `split_words:"true" default:"mysql-mysql"`
	DBPort     int    `split_words:"true" default:"3306"`
	DBUser     string `split_words:"true" default:"root"`
	DBName     string `split_words:"true" default:"users"`
	DBPassword string `split_words:"true" default:"password"`
}

func loadConfig() (*config, error) {
	var cfg config
	err := envconfig.Process("app", &cfg)
	return &cfg, err
}

// ConnectionString returns connection string based on the DBInfo configuration.
func (c *config) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
