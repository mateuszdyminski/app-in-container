package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
)

type app struct {
	config map[string]string
	router *mux.Router
}

func main() {
	app := app{
		config: loadConfig(),
		router: mux.NewRouter(),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for range c {
			log.Println("got a HUP signal")
			app.reload()
		}
	}()

	app.router.HandleFunc("/", app.configHandler)
	app.router.HandleFunc("/config", app.configHandler)
	app.router.HandleFunc("/-/reload", app.reloadHandler).Methods(http.MethodPost)

	// start server
	http.ListenAndServe(":8080", app.router)
}

func (a *app) configHandler(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(a.config)
	if err != nil {
		w.Write([]byte("can't marshal configuration"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (a *app) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(a.config["message"]))
}

func (a *app) reloadHandler(w http.ResponseWriter, r *http.Request) {
	a.reload()
}

func (a *app) reload() {
	a.config = loadConfig()
}

func loadConfig() map[string]string {
	configPath := os.Getenv("CONFIG_PATH")
	log.Printf("loading configuration from: %s\n", configPath)

	files, _ := ioutil.ReadDir(configPath)
	config := make(map[string]string)
	for _, file := range files {
		filename := path.Join(configPath, file.Name())

		if !strings.HasPrefix(file.Name(), ".") {
			value, _ := ioutil.ReadFile(filename)
			config[file.Name()] = string(value)
		}
	}

	return config
}
