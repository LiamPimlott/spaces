package spaces

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

// NewCreateSpaceHandler returns an http handler for creating spaces
func NewCreateSpaceHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spaceReq := &Space{}

		utils.Decode(w, r, spaceReq)

		ok, err := spaceReq.Valid()
		if !ok || err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, err.Error())
			return
		}

		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
		if !ok {
			log.Println("error getting claims")
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		spaceReq.OwnerID = claims.ID

		space, err := s.Create(*spaceReq)
		if err != nil {
			log.Println("error creating space")
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, space)
	}
}

// NewGetSpaceByIDHandler returns an http handler for getting spaces by id
func NewGetSpaceByIDHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rtPrms := mux.Vars(r)

		idStrng, ok := rtPrms["id"]
		if !ok {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "missing space id in url")
			return
		}

		id, err := strconv.Atoi(idStrng)
		if err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "invalid space id in url")
			return
		}

		space, err := s.GetByID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondError(w, errs.ErrNotFound.Code, errs.ErrNotFound.Msg, "")
				return
			}
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
		if !ok || claims.ID != space.OwnerID {
			space = Space{
				ID:          space.ID,
				OwnerID:     space.OwnerID,
				Title:       space.Title,
				Description: space.Description,
				City:        space.City,
				Province:    space.Province,
				Country:     space.Country,
			}
		}

		utils.Respond(w, space)
	}
}
