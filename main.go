package main

import (
	"net/http"
	"time"

	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
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

var log *zap.Logger
var casbinEnforcer *casbin.Enforcer
var tokens map[string]string
var passwords map[string]string

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/generate", GenerateHashHandler)
	r.HandleFunc("/auth", AuthenticationHandler).Methods("GET")
	r.Handle("/{resource}/{name}", basicAuthMiddleware(http.HandlerFunc(ResourceHandler))).Methods("GET")
	http.Handle("/", r)

	httpAddress := viper.GetString("address") + ":" + viper.GetString("port")

	srv := &http.Server{
		Handler:      r,
		Addr:         httpAddress,
		WriteTimeout: writeTimeout * time.Second,
		ReadTimeout:  readTimeout * time.Second,
	}

	log.Info("Starting server", zap.String("address", httpAddress))
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server", zap.Error(err))
	}

}
