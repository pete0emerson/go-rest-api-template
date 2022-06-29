package main

import (
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	data["data"] = "Here is some data you have access to"
	writeData(w, data, http.StatusOK)

}

func GenerateHashHandler(w http.ResponseWriter, r *http.Request) {
	uuid := getUUID()
	data := map[string]string{}

	log.Info("Request",
		zap.String("uuid", uuid),
		zap.String("uri", r.RequestURI),
		zap.String("method", r.Method),
	)

	username, password, ok := r.BasicAuth()
	if !ok {
		log.Error("No basic auth")
		data["error"] = "No credentials provided"
		writeData(w, data, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password", zap.Error(err))
	}
	data["hash"] = string(hashedPassword)
	log.Info("Generated hash from password",
		zap.String("uuid", uuid),
		zap.String("user", username),
		zap.String("hash", string(hashedPassword)),
	)
	// We'll store the first password generated for the user
	if len(passwords) == 0 {
		log.Info("Stored first password hash for user", zap.String("user", username), zap.String("hash", string(hashedPassword)))
		passwords[username] = string(hashedPassword)
	}
	writeData(w, data, http.StatusOK)

}

func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	uuid := getUUID()

	log.Info("Request",
		zap.String("uuid", uuid),
		zap.String("uri", r.RequestURI),
		zap.String("method", r.Method),
	)

	user, pass, ok := r.BasicAuth()
	if !ok || !checkUsernameAndPassword(uuid, user, pass) {
		data := map[string]string{}
		writeData(w, data, http.StatusUnauthorized)
		return
	}

	data := map[string]string{}
	data["token"] = generateSecureToken(tokenLength)
	tokens[user] = data["token"]
	log.Info("Authentication granted",
		zap.String("uuid", uuid),
		zap.String("user", user),
		zap.String("token", data["token"]),
	)
	writeData(w, data, http.StatusOK)

}
