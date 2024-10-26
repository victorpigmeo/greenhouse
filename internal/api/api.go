package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/victorpigmeo/greenhouse/models"
)

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index)
	mux.HandleFunc("POST /api/auth", auth)
	mux.HandleFunc("GET /api/dht", read_dht)

	http.ListenAndServe(":8080", mux)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	slog.Info(fmt.Sprintf(`%s %s`, r.Method, string(r.URL.Path)))

	http.ServeFile(w, r, "./static/index.html")
}

func auth(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/auth" {
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

	authRequest := models.AuthRequest{}
	err = json.NewDecoder(r.Body).Decode(&authRequest)

	if err != nil {
		slog.Error("Error decoding request body")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error decoding request body"))
		return
	}

	if authRequest.Password == config.AdminPassword {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func read_dht(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/dht" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dhtOutput, err := exec.Command("python3", "dht11-raspberry-pi.py").Output()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error fetching DHT data"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dhtOutput))

}
