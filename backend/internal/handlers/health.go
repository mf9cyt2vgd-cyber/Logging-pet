package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RootResponse struct {
	Service   string
	Endpoints struct {
		Status  string
		Docs    string
		Metrics string
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := RootResponse{
		Service: "logs-backend", Endpoints: struct {
			Status  string
			Docs    string
			Metrics string
		}{Status: "/status/{code}", Docs: "/docs", Metrics: "/metrics"},
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		slog.Error("error encoding response")
	}
}
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	code, err := strconv.Atoi(parts[2])
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid status code"})
		return
	}
	seconds := r.URL.Query().Get("seconds_sleep")
	if seconds != "" {
		sec, _ := strconv.Atoi(seconds)
		time.Sleep(time.Duration(sec) * time.Second)
	}
	slog.Info("Hello from Go!", "status code", code, "seconds sleep", seconds)
	if code == 200 {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{
			"data": "Hello",
		})
		if err != nil {
			slog.Error("error encoding response")
		}
		return
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"error": "an error occurred"})
	if err != nil {
		slog.Error("error encoding response")
	}
	return
}
