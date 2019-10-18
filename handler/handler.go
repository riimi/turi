package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

func Home(buildTime, commit, release string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := struct {
			BuildTime string `json:"buildTime"`
			Commit    string `json:"commit"`
			Release   string `json:"release"`
		}{
			BuildTime: buildTime, Commit: commit, Release: release,
		}

		b, err := json.Marshal(info)
		if err != nil {
			log.Printf("Failed to encode info data: %v", err)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func Health(isHealthy *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isHealthy == nil || !isHealthy.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func Healthy(isHealthy *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isHealthy.Store(true)
	}
}

func Unhealthy(isHealthy *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isHealthy.Store(false)
	}
}

func Readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func marshalJson(x interface{}) string {
	if x == nil {
		return ""
	}
	b, _ := json.Marshal(x)
	return string(b)
}

func bindJSON(r *http.Request, x interface{}) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(x)
}
