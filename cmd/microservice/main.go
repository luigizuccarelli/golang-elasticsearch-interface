// +build real

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/connectors"
	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/handlers"
	"lmzsoftware.com/luigizuccarelli/golang-elasticsearch-interface/pkg/validator"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

var (
	logger       *simple.Logger
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
		// use this for cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Language")
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

// startHttpServer - private function - http start
func startHttpServer(con connectors.Clients) *http.Server {
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	r := mux.NewRouter()

	// setup prometheus endpoint
	r.Use(prometheusMiddleware)
	r.Path("/metrics").Handler(promhttp.Handler())

	// used for customized search
	r.HandleFunc("/api/v1/search/customers", func(w http.ResponseWriter, req *http.Request) {
		handlers.SearchHandler(w, req, con)
	}).Methods("POST", "GET", "OPTIONS")

	// read and update/insert single document endpoint
	r.HandleFunc("/api/v1/customers", func(w http.ResponseWriter, req *http.Request) {
		handlers.DocumentHandler(w, req, con)
	}).Methods("PUT", "DELETE", "OPTIONS")

	r.HandleFunc("/api/v2/sys/info/isalive", handlers.IsAlive).Methods("GET")

	http.Handle("/", r)

	if err := srv.ListenAndServe(); err != nil {
		con.Error("Httpserver: ListenAndServe() error: " + err.Error())
		return nil
	}

	return srv
}

// main - needs no explanation :)
func main() {
	var logger *simple.Logger

	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}

	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}

	conn := connectors.NewClientConnectors(logger)

	logger.Info("Starting server on port " + os.Getenv("SERVER_PORT"))
	srv := startHttpServer(conn)
	logger.Info(fmt.Sprintf("Started server %v", srv))
}
