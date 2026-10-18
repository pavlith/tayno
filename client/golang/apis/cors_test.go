package apis

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestSignupCORSPreflight(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:8082")
	s := NewServer(registrationApp{}, fstest.MapFS{"static/index.html": &fstest.MapFile{Data: []byte("test")}})

	for _, tc := range []struct {
		origin string
		allow  bool
	}{
		{"http://localhost:8082", true},
		{"https://untrusted.example", false},
	} {
		t.Run(tc.origin, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodOptions, "/api/auth/signup", nil)
			req.Header.Set("Origin", tc.origin)
			req.Header.Set("Access-Control-Request-Method", "POST")
			req.Header.Set("Access-Control-Request-Headers", "content-type,x-myst-account")
			res := httptest.NewRecorder()
			s.ServeHTTP(res, req)
			if got := res.Header().Get("Access-Control-Allow-Origin"); (got == tc.origin) != tc.allow {
				t.Fatalf("origin %q: Access-Control-Allow-Origin = %q", tc.origin, got)
			}
			if tc.allow {
				if got := res.Header().Get("Access-Control-Allow-Methods"); got == "" {
					t.Fatal("preflight missing allowed methods")
				}
				if got := res.Header().Get("Access-Control-Allow-Headers"); got == "" {
					t.Fatal("preflight missing allowed headers")
				}
			}
		})
	}
}
