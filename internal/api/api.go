package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/victorpigmeo/greenhouse/models"
)

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index)
	mux.HandleFunc("PUT /unlock", unlock)
	mux.HandleFunc("GET /secure", secure)
	mux.HandleFunc("GET /api/dht", readDht)
	mux.HandleFunc("PUT /api/gpio/{action}/{pin}", gpio)

	http.ListenAndServe(":8080", mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	slog.Info(fmt.Sprintf(`%s %s`, r.Method, string(r.URL.Path)))

	http.ServeFile(w, r, "./static/auth.html")
}

func unlock(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/unlock" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	slog.Info(fmt.Sprintf(`%s %s`, r.Method, string(r.URL.Path)))

	file, err := os.Open("./config.json")
	defer file.Close()

	if err != nil {
		slog.Error("Error trying to open config file")
		fmt.Println(err)
	}

	config := models.Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		slog.Error("Error trying to decode config file")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error trying to open config file"))
		return
	}

	password := r.Header.Get("Authentication")

	if password == "" {
		slog.Error("Authentication header is empty")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Authentication header is empty"))
		return
	}

	if password == config.UnlockPassword {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"unlocked\": true}"))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func secure(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/secure" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	slog.Info(fmt.Sprintf(`%s %s`, r.Method, string(r.URL.Path)))

	http.ServeFile(w, r, "./static/index.html")
}

func readDht(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/dht" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	slog.Info(fmt.Sprintf(`%s %s`, r.Method, string(r.URL.Path)))

	res, err := http.Get("http://192.168.18.26:8080/api/dht")

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	resBody, err := io.ReadAll(res.Body)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resBody))

}

func gpio(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.URL.Path, "/api/gpio/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, err := os.Open("./config.json")
	defer file.Close()

	if err != nil {
		slog.Error("Error trying to open config file")
		fmt.Println(err)
	}

	config := models.Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		slog.Error("Error trying to decode config file")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error trying to open config file"))
		return
	}

	if r.Header.Get("Authentication") == "" || r.Header.Get("Authentication") != config.AdminPassword {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	action := r.PathValue("action")
	pin, err := strconv.ParseUint(r.PathValue("pin"), 10, 8)

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := http.Get(fmt.Sprintf("http://192.168.18.26:8080/api/gpio/%s/%d", action, pin))

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	resBody, err := io.ReadAll(res.Body)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resBody))

}
