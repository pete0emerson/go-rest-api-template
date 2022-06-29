package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// basicAuthMiddleware makes sure that the user is authenticated and authorized to access the resource
func basicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		uuid := getUUID()

		log.Info("Request",
			zap.String("uuid", uuid),
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
		)

		data := map[string]string{}

		vars := mux.Vars(r)
		name := vars["name"]
		resource := vars["resource"]

		act := "read"
		if r.Method == "POST" {
			act = "create"
		}

		if !hasAuthentication(uuid, name, r.Header.Get("Token")) {
			writeData(w, data, http.StatusForbidden)
			return
		}

		if !hasAuthorization(uuid, name, resource, act) {
			writeData(w, data, http.StatusUnauthorized)
			return
		}

		log.Info("Access granted",
			zap.String("uuid", uuid),
			zap.String("name", name),
			zap.String("resource", resource),
			zap.String("act", act),
			zap.String("access", "granted"),
		)
		handler(w, r)
	}
}

// hasAuthentication checks if the user is authenticated by matching the user's token against what is stored in memory
func hasAuthentication(uuid string, name string, token string) bool {

	startTime := getTime()

	value, okay := tokens[name]
	if !okay {
		log.Info("No token found",
			zap.String("uuid", uuid),
			zap.String("name", name),
			zap.String("submitted_token", token),
			zap.String("access", "denied"),
		)
		endTimer(startTime, "hasAuthorization")
		return false
	}
	if token != value {
		log.Info("Invalid token",
			zap.String("uuid", uuid),
			zap.String("name", name),
			zap.String("submitted_token", token),
			zap.String("stored_token", value),

			zap.String("access", "denied"),
		)
		endTimer(startTime, "hasAuthorization")
		return false
	}
	endTimer(startTime, "hasAuthentication")
	return true
}

// hasAuthorization uses Casbin to verify that the user is authorized to access the given resource
func hasAuthorization(uuid string, name string, resource string, action string) bool {

	startTime := getTime()

	casbinEnforcer = casbin.NewEnforcer(viper.GetString("auth-model"), viper.GetString("auth-policy"))

	if !casbinEnforcer.Enforce(name, resource, action) {
		log.Info("Unauthorized",
			zap.String("uuid", uuid),
			zap.String("name", name),
			zap.String("resource", resource),
			zap.String("act", action),
			zap.String("access", "denied"),
		)
		endTimer(startTime, "hasAuthorization")
		return false
	}
	endTimer(startTime, "hasAuthorization")
	return true
}

// generateSecureToken generates a random secure token for the user
func generateSecureToken(length int) string {
	startTime := getTime()
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	endTimer(startTime, "generateSecureToken")
	return hex.EncodeToString(b)
}

// checkUsernameAndPassword uses bcrypt to check if the username and password stored in memory are correct
func checkUsernameAndPassword(uuid, username, password string) bool {

	startTime := getTime()

	// Make sure the user exists
	if _, okay := passwords[username]; !okay {
		log.Info("Invalid username",
			zap.String("uuid", uuid),
			zap.String("username", username),
		)
		endTimer(startTime, "checkUsernameAndPassword")
		return false
	}

	// Hash the password and compare to what is stored
	err := bcrypt.CompareHashAndPassword([]byte(passwords[username]), []byte(password))
	if err != nil {
		log.Info("Invalid password",
			zap.String("uuid", uuid),
			zap.String("username", username),
		)
		endTimer(startTime, "checkUsernameAndPassword")
		return false
	}

	log.Info("Valid username and password",
		zap.String("uuid", uuid),
		zap.String("username", username),
	)
	endTimer(startTime, "checkUsernameAndPassword")
	return true
}
