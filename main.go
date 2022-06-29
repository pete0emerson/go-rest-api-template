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

type Response struct {
	Code   string `json:"code"`
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Status string `json:"status"`
	Token  string `json:"token"`
}

var log *zap.Logger
var casbinEnforcer *casbin.Enforcer
var tokens map[string]string
var passwords map[string]string

func main() {

	// The router definitions and mapping to the handlers
	r := mux.NewRouter()
	r.HandleFunc("/generate", GenerateHashHandler)
	r.HandleFunc("/auth", AuthenticationHandler).Methods("GET")
	r.Handle("/{resource}/{name}", basicAuthMiddleware(http.HandlerFunc(ResourceHandler))).Methods("GET")
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
