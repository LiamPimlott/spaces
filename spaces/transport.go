package spaces

import (
	"log"
	"net/http"

	"github.com/gorilla/context"

	"github.com/LiamPimlott/spaces/lib"
)

// NewCreateSpaceHandler returns an http handler for creating spaces
func NewCreateSpaceHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		spcReq := &Space{}

		utils.Decode(w, r, spcReq)

		ok, err := spcReq.Valid()
		if !ok || err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
		if !ok {
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		spcReq.OwnerID = claims.ID

		spc, err := s.Create(*spcReq)
		if err != nil {
			log.Println("error creating space.")
			w.WriteHeader(500)
			w.Write([]byte("Internal server error"))
			return
		}

		utils.Respond(w, spc)
	}
}
