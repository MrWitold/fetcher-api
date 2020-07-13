package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/MrWitold/fetcher-api/api/models"
	"github.com/MrWitold/fetcher-api/api/responses"
)

// ShowHistory controls flow for history endpoint
func (s *Server) ShowHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	h := models.History{}

	history, err := h.FindHistoryByID(s.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, history)
}
