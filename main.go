package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const defaultPort = 8000
const defaultAddress = "127.0.0.1"
const writeTimeout = 15
const readTimeout = 15
const envPrefix = "server"
const configFileName = "server"

// Paths should be comma separated
const configPaths = ".,./config"
const tokenLength = 20

type Response struct {
	Code         int    `json:"code,omitempty"`
	Data         string `json:"data,omitempty"`
	Hash         string `json:"hash,omitempty"`
	Message      string `json:"message,omitempty"`
	Status       string `json:"status,omitempty"`
	Token        string `json:"token,omitempty"`
	BuildVersion string `json:"version,omitempty"`
	BuildDate    string `json:"build-date,omitempty"`
}

var log *zap.Logger
var casbinEnforcer *casbin.Enforcer
var tokens map[string]string
var passwords map[string]string
var buildVersion string
var buildDate string

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

func main() {

	// The router definitions and mapping to the handlers
	r := mux.NewRouter()

	r.Use(prometheusMiddleware)

	r.HandleFunc("/generate", GenerateHashHandler)
	r.HandleFunc("/auth", AuthenticationHandler).Methods("GET")
	r.Handle("/{resource}/{name}", basicAuthMiddleware(http.HandlerFunc(ResourceHandler))).Methods("GET")
	r.HandleFunc("/redis", RedisHandler)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/version", VersionHandler)
	http.Handle("/", r)

	// Configure the server parameters
	httpAddress := viper.GetString("address") + ":" + viper.GetString("port")

	srv := &http.Server{
		Handler:      r,
		Addr:         httpAddress,
		WriteTimeout: writeTimeout * time.Second,
		ReadTimeout:  readTimeout * time.Second,
	}

	// Fire up the server
	log.Info("Starting server", zap.String("address", httpAddress))
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server", zap.Error(err))
	}

}
