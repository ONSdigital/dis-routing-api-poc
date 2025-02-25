package api

import (
	"context"

	"github.com/ONSdigital/dis-routing-api-poc/store"
	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
type API struct {
	Router *mux.Router
	Store  *store.DataStore
}

// Setup function sets up the API and returns an API
func Setup(ctx context.Context, r *mux.Router, store *store.DataStore) *API {
	api := &API{
		Router: r,
		Store:  store,
	}

	api.Router.HandleFunc("/api/v1/routes", GetAllRoutesHandler(store)).Methods("GET")
	api.Router.HandleFunc("/api/v1/routes/create", CreateRouteHandler(store)).Methods("POST")
	api.Router.HandleFunc("/api/v1/routes/{id}", UpdateRouteHandler(store)).Methods("PUT")
	api.Router.HandleFunc("/api/v1/routes/{id}", GetRouteHandler(store)).Methods("GET")
	api.Router.HandleFunc("/api/v1/routes/delete/{id}", DeleteRouteHandler(store)).Methods("DELETE")

	api.Router.HandleFunc("/api/v1/redirects", GetAllRedirectsHandler(store)).Methods("GET")
	api.Router.HandleFunc("/api/v1/redirects/create", CreateRedirectHandler(store)).Methods("POST")
	api.Router.HandleFunc("/api/v1/redirects/{id}", UpdateRedirectHandler(store)).Methods("PUT")
	api.Router.HandleFunc("/api/v1/redirects/{id}", GetRedirectHandler(store)).Methods("GET")
	api.Router.HandleFunc("/api/v1/redirects/delete/{id}", DeleteRedirectHandler(store)).Methods("DELETE")

	return api
}
