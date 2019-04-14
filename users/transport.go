package users

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/LiamPimlott/spaces/lib/errs"
	"github.com/LiamPimlott/spaces/lib/utils"
)

// NewCreateUserHandler returns an http handler for creating users
func NewCreateUserHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usrReq := &User{}

		utils.Decode(w, r, usrReq)

		ok, err := usrReq.Valid()
		if !ok || err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, err.Error())
			return
		}

		usr, err := s.Create(*usrReq)
		if err != nil {
			log.Println("error creating user")
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, usr)
	}
}

// NewLoginHandler returns an http handler for logging in users
func NewLoginHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := &User{}

		utils.Decode(w, r, body)

		ok := body.ValidLogin()
		if !ok {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "")
			return
		}

		tkn, err := s.Login(*body)
		if err != nil {
			log.Printf("error logging in user: %s\n", err)
			if err == sql.ErrNoRows {
				utils.RespondError(w, errs.ErrNotFound.Code, errs.ErrNotFound.Msg, "")
				return
			}
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, tkn)
	}
}

// NewGetUserByIDHandler returns an http handler for getting users by id
func NewGetUserByIDHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rtPrms := mux.Vars(r)

		idStrng, ok := rtPrms["id"]
		if !ok {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "missing user id in url")
			return
		}

		id, err := strconv.Atoi(idStrng)
		if err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "invalid user id in url")
			return
		}

		usr, err := s.GetByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondError(w, errs.ErrNotFound.Code, errs.ErrNotFound.Msg, "")
				return
			}
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
		if !ok {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		if claims.ID != usr.ID {
			usr = User{
				FirstName: usr.FirstName,
				LastName:  usr.LastName,
			}
		}

		utils.Respond(w, usr)
	}
}
