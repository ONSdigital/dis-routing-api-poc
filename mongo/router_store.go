package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/ONSdigital/dis-routing-api-poc/config"
	"github.com/ONSdigital/dis-routing-api-poc/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	mongohealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	mongodriver.MongoDriverConfig

	Connection   *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient
}

// NewDBConnection creates a new Mongo object encapsulating a connection to the mongo server/cluster with the given configuration,
// and a health client to check the health of the mongo server/cluster
func NewDBConnection(_ context.Context, cfg config.MongoConfig) (m *Mongo, err error) {
	m = &Mongo{MongoDriverConfig: cfg}
	m.Connection, err = mongodriver.Open(&m.MongoDriverConfig)
	if err != nil {
		return nil, err
	}

	databaseCollectionBuilder := map[mongohealth.Database][]mongohealth.Collection{
		mongohealth.Database(m.Database): {
			mongohealth.Collection(m.ActualCollectionName(config.RoutesCollection)),
			mongohealth.Collection(m.ActualCollectionName(config.RedirectsCollection)),
		},
	}
	m.healthClient = mongohealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

	return m, nil
}

// Close closes the mongo session and returns any error
// It is an error to call m.Close if m.Init() returned an error, and there is no open connection
func (m *Mongo) Close(ctx context.Context) error {
	return m.Connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// ValidateRoute performs validation on a Route object, including checks for path conflicts
func (m *Mongo) ValidateRoute(ctx context.Context, route *models.Route) error {
	if route.Path == "" {
		return errors.New("route path cannot be empty")
	}
	if len(route.Domains) == 0 {
		return errors.New("at least one domain is required")
	}

	// Prevent overlapping routes
	var existingRoute models.Route
	err := m.Connection.Collection(config.RoutesCollection).FindOne(ctx, bson.M{"path": route.Path}, &existingRoute)
	if err != nil && err != mongo.ErrNoDocuments {
		// An unexpected error occurred
		return fmt.Errorf("error checking existing routes: %w", err)
	}
	if err == nil {
		// If no error, it means a route with the same path exists
		return errors.New("overlapping route for the same path already exists")
	}

	return nil
}

// ValidateRedirect performs validation on a Redirect object, including checks for circular redirects
func (m *Mongo) ValidateRedirect(ctx context.Context, redirect *models.Redirect) error {
	if redirect.From == "" || redirect.To == "" {
		return errors.New("redirect 'from' and 'to' paths cannot be empty")
	}
	if redirect.StatusCode != 307 && redirect.StatusCode != 308 {
		return errors.New("redirect status code must be 307 or 308")
	}

	// Prevent circular redirects by checking if 'To' already redirects back to 'From'
	var existingRedirect models.Redirect
	err := m.Connection.Collection(config.RedirectsCollection).FindOne(ctx, bson.M{"from": redirect.To, "to": redirect.From}, &existingRedirect)
	if err == nil {
		return errors.New("circular redirect detected")
	}

	// Prevent overlapping redirects
	var overlappingRedirect models.Redirect
	err = m.Connection.Collection(config.RedirectsCollection).FindOne(ctx, bson.M{"from": redirect.From}, &overlappingRedirect)
	if err == nil {
		return errors.New("overlapping redirect from the same path already exists")
	}

	return nil
}

// GetAllRoutes retrieves all routes with optional filters using dp-mongodb's Find method
func (m *Mongo) GetAllRoutes(ctx context.Context, filter map[string]interface{}) (*[]models.Route, error) {
	var routes []models.Route
	total, err := m.Connection.Collection(config.RoutesCollection).Find(ctx, filter, &routes)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &[]models.Route{}, nil
	}

	return &routes, nil
}

// GetRoute retrieves a route by its ID
func (m *Mongo) GetRoute(ctx context.Context, id string) (*models.Route, error) {
	var route models.Route
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = m.Connection.Collection(config.RoutesCollection).FindOne(ctx, bson.M{"_id": objID}, &route)
	if err != nil {
		return nil, err
	}
	return &route, nil
}

// CreateRoute inserts a new route
func (m *Mongo) CreateRoute(ctx context.Context, route *models.Route) error {
	route.CreatedAt = time.Now()
	route.UpdatedAt = time.Now()
	_, err := m.Connection.Collection(config.RoutesCollection).InsertOne(ctx, route)
	return err
}

// UpdateRoute modifies an existing route
func (m *Mongo) UpdateRoute(ctx context.Context, id string, update map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	// Check if "path" is being updated
	if newPath, ok := update["path"].(string); ok {
		var existingRoute models.Route
		err := m.Connection.Collection(config.RoutesCollection).FindOne(ctx, bson.M{"path": newPath}, &existingRoute)
		if err == nil && existingRoute.ID != id {
			// If a different route already exists with this path, reject update
			return errors.New("another route with this path already exists")
		} else if err != nil && err != mongo.ErrNoDocuments {
			// If there's an unexpected DB error, return it
			return fmt.Errorf("error checking existing routes: %w", err)
		}
	}

	update["updated_at"] = time.Now()
	_, err = m.Connection.Collection(config.RoutesCollection).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	return err
}

// DeleteRoute removes a route
func (m *Mongo) DeleteRoute(ctx context.Context, id string) error {
	// Convert the string id to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}
	_, err = m.Connection.Collection(config.RoutesCollection).DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// GetAllRedirects retrieves all redirects with optional filters using dp-mongodb's Find method
func (m *Mongo) GetAllRedirects(ctx context.Context, filter map[string]interface{}) (*[]models.Redirect, error) {
	var redirects []models.Redirect
	total, err := m.Connection.Collection(config.RedirectsCollection).Find(ctx, filter, &redirects)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &[]models.Redirect{}, nil
	}

	return &redirects, nil
}

// GetRedirect retrieves a redirect by its ID
func (m *Mongo) GetRedirect(ctx context.Context, id string) (*models.Redirect, error) {
	var redirect models.Redirect
	err := m.Connection.Collection(config.RedirectsCollection).FindOne(ctx, bson.M{"_id": id}, &redirect)
	if err != nil {
		return nil, err
	}
	return &redirect, nil
}

// CreateRedirect inserts a new redirect
func (m *Mongo) CreateRedirect(ctx context.Context, redirect *models.Redirect) error {
	redirect.CreatedAt = time.Now()
	redirect.UpdatedAt = time.Now()
	_, err := m.Connection.Collection(config.RedirectsCollection).InsertOne(ctx, redirect)
	return err
}

// UpdateRedirect modifies an existing redirect
func (m *Mongo) UpdateRedirect(ctx context.Context, id string, update map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}
	update["updated_at"] = time.Now()
	_, err = m.Connection.Collection(config.RedirectsCollection).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	return err
}

// DeleteRedirect removes a redirect
func (m *Mongo) DeleteRedirect(ctx context.Context, id string) error {
	// Convert the string id to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	_, err = m.Connection.Collection(config.RedirectsCollection).DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
