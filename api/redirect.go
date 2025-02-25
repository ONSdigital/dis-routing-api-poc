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

// Redirect Handlers
func CreateRedirectHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var redirect models.Redirect
		if err := json.NewDecoder(r.Body).Decode(&redirect); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if err := store.Backend.ValidateRedirect(context.Background(), &redirect); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := store.Backend.CreateRedirect(context.Background(), &redirect); err != nil {
			http.Error(w, "Failed to create redirect", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(redirect); err != nil {
			log.Printf("Failed to encode redirect response: %v", err)
		}
	}
}

func GetRedirectHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "Missing redirect ID", http.StatusBadRequest)
			return
		}
		redirect, err := store.Backend.GetRedirect(context.Background(), id)
		if err != nil {
			http.Error(w, "Redirect not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(redirect); err != nil {
			log.Printf("Failed to encode redirect response: %v", err)
		}
	}
}

func GetAllRedirectsHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirects, err := store.Backend.GetAllRedirects(context.Background(), map[string]interface{}{})
		if err != nil {
			http.Error(w, "Failed to fetch redirects", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(redirects); err != nil {
			http.Error(w, "Failed to encode redirects", http.StatusInternalServerError)
			return
		}
	}
}

func DeleteRedirectHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "Missing redirect ID", http.StatusBadRequest)
			return
		}
		if err := store.Backend.DeleteRedirect(context.Background(), id); err != nil {
			http.Error(w, "Failed to delete redirect", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func UpdateRedirectHandler(store *store.DataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "Missing redirect ID", http.StatusBadRequest)
			return
		}
		var update map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if err := store.Backend.UpdateRedirect(context.Background(), id, update); err != nil {
			http.Error(w, "Failed to update redirect", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
