package main

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func ResourceHandler(w http.ResponseWriter, r *http.Request) {

	startTime := getTime()

	uuid := getUUID()

	log.Info("ResourceHandler Request",
		zap.String("uuid", uuid),
		zap.String("uri", r.RequestURI),
		zap.String("method", r.Method),
	)

	data := Response{}
	data.Data = "Here is some data you have access to"
	writeData(w, data, http.StatusOK)

	endTimer(startTime, "ResourceHandler")

}

// GenerateHashHandler generates a hash of the password provided. This would never be used in production as coded below.
func GenerateHashHandler(w http.ResponseWriter, r *http.Request) {

	startTime := getTime()

	uuid := getUUID()
	data := Response{}

	log.Info("GenerateHashHandler Request",
		zap.String("uuid", uuid),
		zap.String("uri", r.RequestURI),
		zap.String("method", r.Method),
	)
	// Get the username and password from the request
	username, password, ok := r.BasicAuth()
	if !ok {
		log.Info("No basic auth")
		data.Status = "No credentials provided"
		writeData(w, data, http.StatusBadRequest)
		endTimer(startTime, "GenerateHashHandler")
		return
	}

	// Generate a hash of the password. We never want to store raw passwords.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error hashing password", zap.Error(err))
	}
	data.Hash = string(hashedPassword)
	log.Info("Generated hash from password",
		zap.String("uuid", uuid),
		zap.String("user", username),
		zap.String("hash", string(hashedPassword)),
	)

	// This allows anyone to store a password. Obviously not secure, this is only for demo purposes.
	log.Info("Stored first password hash for user", zap.String("user", username), zap.String("hash", string(hashedPassword)))
	passwords[username] = string(hashedPassword)
	writeData(w, data, http.StatusOK)

	endTimer(startTime, "GenerateHashHandler")

}

// AuthenticateHandler authenticates the user with the provided username and password.
func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {

	startTime := getTime()

	uuid := getUUID()

	log.Info("Request",
		zap.String("uuid", uuid),
		zap.String("uri", r.RequestURI),
		zap.String("method", r.Method),
	)

	// Get the username and password from the request
	user, pass, ok := r.BasicAuth()
	if !ok || !checkUsernameAndPassword(uuid, user, pass) {
		data := Response{}
		writeData(w, data, http.StatusUnauthorized)
		elapsedTime := time.Since(startTime)
		log.Info("GenerateHashHandler", zap.Duration("elapsed", elapsedTime))
		return
	}

	// Generate a token for the user and return it
	data := Response{}
	data.Token = generateSecureToken(tokenLength)
	tokens[user] = data.Token
	log.Info("Authentication granted",
		zap.String("uuid", uuid),
		zap.String("user", user),
		zap.String("token", data.Token),
	)
	writeData(w, data, http.StatusOK)

	endTimer(startTime, "AuthenticationHandler")

}
