package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// AppRest is REST controller used to handle requests about the application.
type AppRest struct {
	hostname  string
	startedAt time.Time
	db        *sql.DB
	mu        sync.Mutex
	buildInfo BuildInfo
	healthy   bool
}

// NewUserRest constructs new UserRest controller used to handler request about the user entities.
func NewAppRest(db *sql.DB) (*AppRest, error) {
	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &AppRest{
		db:        db,
		hostname:  host,
		startedAt: time.Now().UTC(),
		buildInfo: buildInfo(),
		healthy:   true,
	}, nil
}

// Unhealthy sets server health to false - used in gracefull shutdown mode
func (r *AppRest) Unhealthy() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.healthy = false
}

func (r *AppRest) Ready(w http.ResponseWriter, req *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.healthy {
		w.Write([]byte("ready"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("not-ready"))
	}
}

// Health handler is responsible for serving resonse with current health status of service. Healthz concept is used to leverage the regular health status pattern.
func (r *AppRest) Health(w http.ResponseWriter, req *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.healthy {
		dbstatus := "ok"
		err := r.db.Ping()
		if err != nil {
			dbstatus = "down"
		}

		resp := HealthStatus{
			Hostname:  r.hostname,
			StartedAt: r.startedAt.Format("2006-01-02_15:04:05"),
			Uptime:    time.Now().UTC().Sub(r.startedAt).String(),
			Build:     r.buildInfo,
			DBStatus:  dbstatus,
		}

		WriteJSON(w, resp)
	} else {
		WriteErr(w, fmt.Errorf("Server in graceful  shutdown mode"), http.StatusInternalServerError)
	}
}

// Err handler is dummy simulator of error which occurs in out service.
func (r *AppRest) Err(w http.ResponseWriter, req *http.Request) {
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		WriteErr(w, errors.New("can't read request body"), http.StatusBadRequest)
		return
	}

	WriteErr(w, errors.New(string(bytes)), http.StatusInternalServerError)
}

// HealthStatus holds basic info about the Health status of the application.
type HealthStatus struct {
	Build     BuildInfo `json:"buildInfo"`
	Hostname  string    `json:"hostname"`
	Uptime    string    `json:"uptime"`
	StartedAt string    `json:"startedAt"`
	DBStatus  string    `json:"dbStatus"`
}

// BuildInfo holds basic info about the build based on the git statistics.
type BuildInfo struct {
	Version    string `json:"version"`
	GitVersion string `json:"gitVersion"`
	BuildTime  string `json:"buildTime"`
	LastCommit Commit `json:"lastCommit"`
}

func buildInfo() BuildInfo {
	return BuildInfo{
		Version:    AppVersion,
		GitVersion: GitVersion,
		BuildTime:  BuildTime,
		LastCommit: Commit{
			Author: LastCommitUser,
			ID:     LastCommitHash,
			Time:   LastCommitTime,
		},
	}
}

// Commit holds info about the git commit.
type Commit struct {
	ID     string `json:"id"`
	Time   string `json:"time"`
	Author string `json:"author"`
}

// WriteJSON writes JSON response struct to ResponseWriter.
func WriteJSON(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := w.Write(json); err != nil {
		return err
	}

	return nil
}

// WriteErr writes error to ResponseWriter.
func WriteErr(w http.ResponseWriter, err error, httpCode int) {
	log.Println(err.Error())

	// write error to response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var errMap = map[string]interface{}{
		"httpStatus": httpCode,
		"error":      err.Error(),
	}

	errJSON, err := json.Marshal(errMap)
	if err != nil {
		log.Printf("can't marshal error response. err: %s", err)
	}

	log.Println(string(errJSON))
	w.WriteHeader(httpCode)
	w.Write(errJSON)
}
