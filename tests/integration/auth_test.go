package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"leetcode-spaced-repetition/internal"

	"github.com/gin-gonic/gin"
)

// guardedRoutes are the endpoints that read or write user-scoped data. Every one of them
// must refuse a request the reverse proxy has not authenticated.
var guardedRoutes = []struct {
	method string
	path   string
}{
	{http.MethodGet, "/dashboard"},
	{http.MethodGet, "/problems"},
	{http.MethodGet, "/problems/topics"},
	{http.MethodGet, "/problems/submissions"},
	{http.MethodPost, "/problems/submissions"},
	{http.MethodPost, "/problems/submissions/import"},
	{http.MethodGet, "/problems/1"},
	{http.MethodGet, "/problems/1/submissions"},
}

func TestGuardedRoutesRejectMissingRemoteUser(t *testing.T) {
	for _, route := range guardedRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(route.method, route.path, http.NoBody)
			rawRouter.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected 401 without a %s header, got %d", internal.RemoteUserHeader, w.Code)
			}
		})
	}
}

func TestGuardedRoutesRejectWrongRemoteUser(t *testing.T) {
	for _, route := range guardedRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(route.method, route.path, http.NoBody)
			req.Header.Set(internal.RemoteUserHeader, "someone-else")
			rawRouter.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected 401 for a non-owner user, got %d", w.Code)
			}
		})
	}
}

func TestGuardedRouteAcceptsOwner(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/dashboard", http.NoBody)
	req.Header.Set(internal.RemoteUserHeader, testOwnerUsername)
	rawRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for the owner, got %d: %s", w.Code, w.Body.String())
	}
}

// The middleware must never fall open when it is misconfigured with an empty owner, since a
// request with no Remote-User header would otherwise compare equal to it.
func TestMiddlewareFailsClosedOnEmptyOwnerUsername(t *testing.T) {
	router := gin.New()
	router.GET(
		"/guarded",
		internal.OwnerOnlyAuthMiddleware("", testOwnerUserID),
		func(c *gin.Context) { c.Status(http.StatusOK) },
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/guarded", http.NoBody)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 when owner username is misconfigured as empty, got %d", w.Code)
	}
}
