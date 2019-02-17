package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/LiamPimlott/spaces/lib"
)

// NewCreateUserHandler returns an http handler for creating users
func NewCreateUserHandler(s UsersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usrReq := &User{}

		utils.Decode(w, r, usrReq)

		usr, err := s.Create(*usrReq)
		if err != nil {
			log.Println("error creating user.")
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
		}

		utils.Respond(w, usr)
	}
}

// NewLoginHandler returns an http handler for logging in users
func NewLoginHandler(s UsersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := &User{}

		utils.Decode(w, r, body)

		tkn, err := s.Login(*body)
		if err != nil {
			log.Printf("error logging in user: %s\n", err)

			if err == sql.ErrNoRows {
				w.WriteHeader(404)
				w.Write([]byte("Not found"))
				return
			}

			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
		}

		utils.Respond(w, tkn)
	}
}

// NewGetUserByIDHandler returns an http handler for getting users by id
func NewGetUserByIDHandler(s UsersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rtPrms := mux.Vars(r)

		idStrng, ok := rtPrms["id"]
		if !ok {
			log.Println("missing id")
			w.WriteHeader(400)
			w.Write([]byte("Bad Request."))
		}

		id, err := strconv.Atoi(idStrng)
		if err != nil {
			log.Println("invalid id")
			w.WriteHeader(400)
			w.Write([]byte("Bad Request."))
		}

		uRsp, err := s.GetById(id)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(404)
				w.Write([]byte("Not found"))
				return
			}
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		w.Header().Set("Content-Set", "application/json; charset=utf-8")

		err = json.NewEncoder(w).Encode(uRsp)
		if err != nil {
			log.Println("error encoding json.")
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
		}
	}
}
