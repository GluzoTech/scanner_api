package scanevent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// These cases all fail before the repository is reached, so a nil *sql.DB is
// enough to exercise routing and validation without a database.
func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router, NewHandler(NewService(NewRepository(nil))))
	return router
}

func post(t *testing.T, router *gin.Engine, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestCreateScanEventValidation(t *testing.T) {
	const validUUID = `"event_id":"3f8a1c92-6b47-4d1e-9a30-7c5e2b4f8d61"`
	const validTime = `"timestamp":"2026-09-10T09:14:27.482Z"`

	tests := []struct {
		name string
		body string
		want int
	}{
		{"malformed json", `{`, http.StatusBadRequest},
		{"missing event_id", `{` + validTime + `,"reference_model":"M","barcode_input":"1"}`, http.StatusBadRequest},
		{"event_id not a uuid", `{"event_id":"nope",` + validTime + `,"reference_model":"M","barcode_input":"1"}`, http.StatusBadRequest},
		{"missing timestamp", `{` + validUUID + `,"reference_model":"M","barcode_input":"1"}`, http.StatusBadRequest},
		{"unparseable timestamp", `{` + validUUID + `,"timestamp":"10-09-2026","reference_model":"M","barcode_input":"1"}`, http.StatusBadRequest},
		{"missing reference_model", `{` + validUUID + `,` + validTime + `,"barcode_input":"1"}`, http.StatusBadRequest},
		{"whitespace reference_model", `{` + validUUID + `,` + validTime + `,"reference_model":"   ","barcode_input":"1"}`, http.StatusBadRequest},
		{"whitespace barcode_input", `{` + validUUID + `,` + validTime + `,"reference_model":"M","barcode_input":"  "}`, http.StatusBadRequest},
		{"reference_model too long", `{` + validUUID + `,` + validTime + `,"reference_model":"` + strings.Repeat("x", 256) + `","barcode_input":"1"}`, http.StatusBadRequest},
	}

	router := newTestRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := post(t, router, "/scan/event", tt.body)
			if rec.Code != tt.want {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tt.want, rec.Body.String())
			}
		})
	}
}

func TestRouteIsMountedAtScanEvent(t *testing.T) {
	router := newTestRouter()

	routes := router.Routes()
	found := false
	for _, r := range routes {
		if r.Method == http.MethodPost && r.Path == "/scan/event" {
			found = true
		}
	}
	if !found {
		t.Fatalf("POST /scan/event not registered; got %v", routes)
	}

	if rec := post(t, router, "/scan/events", `{}`); rec.Code != http.StatusNotFound {
		t.Fatalf("unexpected route /scan/events answered %d", rec.Code)
	}
}
