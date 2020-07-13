package controllers

import (
	"net/http"

	"github.com/MrWitold/fetcher-api/api/middlewares"
)

func (s *Server) initializeRoutes() {
	getR := s.Router.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/api/fetcher", middlewares.SetMiddlewareJSON(s.ShowAllLinks))
	getR.HandleFunc("/api/fetcher/{id:[0-9]+}", middlewares.SetMiddlewareJSON(s.ShowLinkByID))
	getR.HandleFunc("/api/fetcher/{id:[0-9]+}/history", middlewares.SetMiddlewareJSON(s.ShowHistory))

	deleteR := s.Router.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/api/fetcher/{id:[0-9]+}", middlewares.SetMiddlewareJSON(s.DeleteLink))

	postR := s.Router.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/api/fetcher", middlewares.SetMiddlewareJSON(s.CreateOrUpdateLink))
}
