package main

import (
	"net/http"
	"time"

	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
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

func main() {

	// The router definitions and mapping to the handlers
	r := mux.NewRouter()
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
