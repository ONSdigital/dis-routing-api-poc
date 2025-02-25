package store

import (
	"context"

	"github.com/ONSdigital/dis-routing-api-poc/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

// DataStore provides an interface used to store, retrieve, remove, or update routes and redirects
// It abstracts the backend implementation (e.g., MongoDB) from the handlers.
type DataStore struct {
	Backend Storer
}

// dataMongoDB represents the required methods to access data from mongoDB
type dataMongoDB interface {
	GetRoute(ctx context.Context, id string) (*models.Route, error)
	GetAllRoutes(ctx context.Context, filter map[string]interface{}) (*[]models.Route, error)
	CreateRoute(ctx context.Context, route *models.Route) error
	UpdateRoute(ctx context.Context, id string, update map[string]interface{}) error
	DeleteRoute(ctx context.Context, id string) error
	ValidateRoute(ctx context.Context, route *models.Route) error

	GetRedirect(ctx context.Context, id string) (*models.Redirect, error)
	GetAllRedirects(ctx context.Context, filter map[string]interface{}) (*[]models.Redirect, error)
	CreateRedirect(ctx context.Context, redirect *models.Redirect) error
	UpdateRedirect(ctx context.Context, id string, update map[string]interface{}) error
	DeleteRedirect(ctx context.Context, id string) error
	ValidateRedirect(ctx context.Context, redirect *models.Redirect) error

	Checker(ctx context.Context, state *healthcheck.CheckState) error
	Close(ctx context.Context) error
}

// MongoDB represents all the required methods from mongo DB
type MongoDB interface {
	dataMongoDB
	Close(context.Context) error
	Checker(context.Context, *healthcheck.CheckState) error
}

// Storer represents basic data access via Get, Remove and Upsert methods, abstracting it from mongoDB
type Storer interface {
	dataMongoDB
}
