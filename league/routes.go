package league

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/league/create", h.createLeagues).Methods(http.MethodPost)
}
