package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strconv"

	"github.com/stianeikeland/go-rpio/v4"
	"github.com/victorpigmeo/greenhouse/models"
)

var originAllowlist = []string{
	"http://127.0.0.1:8080",
	"http://localhost:8080",
}

func checkCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if slices.Contains(originAllowlist, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}
		next.ServeHTTP(w, r)
	})
}

func Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index)
	mux.HandleFunc("POST /api/auth", auth)
	mux.HandleFunc("GET /api/dht", readDht)
	mux.HandleFunc("PUT /api/gpio/{pin}", gpio)

	http.ListenAndServe(":8080", checkCORS(mux))
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

func readDht(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/dht" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	slog.Info(fmt.Sprintf(`%s %s`, r.Method, string(r.URL.Path)))

	dhtOutput, err := exec.Command("python3", "dht11-raspberry-pi.py").Output()

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(dhtOutput))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dhtOutput))

}

func gpio(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/gpio" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pin, err := strconv.ParseUint(r.PathValue("pin"), 10, 8)

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	rpio.Pin(pin).Toggle()
}
