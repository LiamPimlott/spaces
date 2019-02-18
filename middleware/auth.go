package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"

	"github.com/LiamPimlott/spaces/lib"
)

// Authorized middleware to require a valid jwt token
func Authorized(endpoint func(http.ResponseWriter, *http.Request), secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")

		if len(splitToken) != 2 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid Authorization Header")
			return
		}
		reqToken = splitToken[1]

		if reqToken == "" {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid Authorization Header")
			return
		}

		token, err := jwt.ParseWithClaims(reqToken, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte(secret), nil
		})

		if err != nil {
			log.Printf("error authorizing user: %s\n", err.Error())
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid Authorization Header")
			return
		}

		if !token.Valid {
			w.WriteHeader(403)
			fmt.Fprintf(w, "Not Authorized")
			return
		}

		claims, ok := token.Claims.(*utils.CustomClaims)
		if !ok {
			w.WriteHeader(500)
			fmt.Fprintf(w, "Internal Server Error")

		} else if !token.Valid {
			w.WriteHeader(403)
			fmt.Fprintf(w, "Not Authorized")
			return
		}

		context.Set(r, "claims", claims)

		endpoint(w, r)
	})
}
