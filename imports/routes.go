package imports

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/v1/league/{id}/import", h.importLeagues).Methods(http.MethodPost)
	router.HandleFunc("/import/healthcheck", h.healthcheck).Methods(http.MethodGet)
}

func (h *Handler) healthcheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("healthcheck request received")
}
