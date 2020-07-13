package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/MrWitold/fetcher-api/api/models"
	"github.com/MrWitold/fetcher-api/api/responses"
)

const maxFileSize = 1000000

// ShowAllLinks returns all links form database
func (s *Server) ShowAllLinks(w http.ResponseWriter, r *http.Request) {
	l := models.Link{}

	links, err := l.FindAllLinks(s.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, links)
}

// ShowLinkByID returns one link form database
func (s *Server) ShowLinkByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	l := models.Link{}

	link, err := l.FindLinkByID(s.DB, id)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, err)
		return
	}
	responses.JSON(w, http.StatusOK, link)
}

// DeleteLink delete link with specfied id form database
func (s *Server) DeleteLink(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	l := models.Link{}

	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = l.DeleteLink(s.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	s.removeJob(uid)
	responses.JSON(w, http.StatusNoContent, "")
}

// CreateOrUpdateLink create or update existing link
func (s *Server) CreateOrUpdateLink(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(http.MaxBytesReader(w, r.Body, maxFileSize))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	l := models.Link{}
	err = json.Unmarshal(body, &l)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	link, err := l.CreateOrUpdateLink(s.DB, l.URL, l.Interval)
	if err != nil {
		responses.ERROR(w, http.StatusRequestEntityTooLarge, err)
		return
	}
	s.addJob(link)
	responses.JSON(w, http.StatusOK, link.ID)
}
