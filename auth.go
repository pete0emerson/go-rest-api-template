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

func hasAuthentication(uuid string, name string, token string) bool {
	value, okay := tokens[name]
	if !okay {
		log.Info("No token found",
			zap.String("uuid", uuid),
			zap.String("name", name),
			zap.String("submitted_token", token),
			zap.String("access", "denied"),
		)

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

		return false
	}
	return true
}

func hasAuthorization(uuid string, name string, resource string, action string) bool {
	casbinEnforcer = casbin.NewEnforcer(viper.GetString("auth-model"), viper.GetString("auth-policy"))

	if !casbinEnforcer.Enforce(name, resource, action) {
		log.Info("Unauthorized",
			zap.String("uuid", uuid),
			zap.String("name", name),
			zap.String("resource", resource),
			zap.String("act", action),
			zap.String("access", "denied"),
		)
		return false
	}
	return true
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func checkUsernameAndPassword(uuid, username, password string) bool {
	// Make sure the user exists
	if _, okay := passwords[username]; !okay {
		log.Info("Invalid username",
			zap.String("uuid", uuid),
			zap.String("username", username),
		)
		return false
	}
	// Hash the password and compare to what is stored
	err := bcrypt.CompareHashAndPassword([]byte(passwords[username]), []byte(password))
	if err != nil {
		log.Info("Invalid password",
			zap.String("uuid", uuid),
			zap.String("username", username),
		)
		return false
	}

	log.Info("Valid username and password",
		zap.String("uuid", uuid),
		zap.String("username", username),
	)
	return true
}
