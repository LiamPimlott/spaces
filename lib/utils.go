package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// NotImplementedHandler writes a not implemented response
var NotImplementedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

// Decode decodes a json request body
func Decode(w http.ResponseWriter, r *http.Request, body interface{}) {
	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		log.Println("error decoding json.")
		WriteInternalError(w)
		return
	}
}

// Respond encodes and writes a json response body
func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("error encoding json: %s", err.Error())
		WriteInternalError(w)
	}
}

// Message creates a map containing a status and message
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// WriteInternalError writes an 500 internal error to a ResponseWriter
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("Internal server error"))
}

// GenerateToken generates a signed jwt token string for a user
func GenerateToken(userName, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = userName
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	ts, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("error generating token: %s", err.Error())
		return "", err
	}

	return ts, nil
}