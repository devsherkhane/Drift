package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify all security headers are present
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"),
		"X-Content-Type-Options should be 'nosniff'")

	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"),
		"X-Frame-Options should be 'DENY'")

	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"),
		"X-XSS-Protection should be '1; mode=block'")

	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"),
		"Referrer-Policy should be 'strict-origin-when-cross-origin'")

	assert.Equal(t, "camera=(), microphone=(), geolocation=()", w.Header().Get("Permissions-Policy"),
		"Permissions-Policy should restrict camera, microphone, geolocation")
}
