package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ONSdigital/dis-routing-api-poc/models"
	"github.com/ONSdigital/dis-routing-api-poc/store"
	"github.com/gorilla/mux"
)

func CreateRouteHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var route models.Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if err := store.Backend.ValidateRoute(context.Background(), &route); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := store.Backend.CreateRoute(context.Background(), &route); err != nil {
			http.Error(w, "Failed to create route", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(route); err != nil {
			log.Printf("Failed to encode route response: %v", err)
		}
	}
}

func GetRouteHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "Missing route ID", http.StatusBadRequest)
			return
		}
		route, err := store.Backend.GetRoute(context.Background(), id)
		if err != nil {
			http.Error(w, "Route not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(route); err != nil {
			log.Printf("Failed to encode route response: %v", err)
		}
	}
}

func GetAllRoutesHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes, err := store.Backend.GetAllRoutes(context.Background(), map[string]interface{}{})
		if err != nil {
			http.Error(w, "Failed to fetch routes", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(routes); err != nil {
			log.Printf("Failed to encode routes response: %v", err)
		}
	}
}

func DeleteRouteHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "Missing route ID", http.StatusBadRequest)
			return
		}
		if err := store.Backend.DeleteRoute(context.Background(), id); err != nil {
			http.Error(w, "Failed to delete route", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateRouteHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "Missing route ID", http.StatusBadRequest)
			return
		}
		var update map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if err := store.Backend.UpdateRoute(context.Background(), id, update); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
