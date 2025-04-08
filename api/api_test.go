package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dis-routing-api-poc/store"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := mux.NewRouter()
		ctx := context.Background()
		store := &store.DataStore{}
		api := Setup(ctx, r, store)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/api/v1/routes", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/api/v1/redirects", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, http.NoBody)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
