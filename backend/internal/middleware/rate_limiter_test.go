package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_AllowsNormalTraffic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	// Allow 5 requests/sec with burst of 5
	r.Use(RateLimiter(5, 5))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// First 5 requests should succeed (burst capacity)
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}
}

func TestRateLimiter_BlocksExcessiveTraffic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	// Allow 1 request/sec with burst of 2
	r.Use(RateLimiter(1, 2))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// First 2 requests fill the burst
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:5678"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// 3rd request should be rate-limited
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:5678"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "Too many requests")
}

func TestRateLimiter_IndependentPerIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RateLimiter(1, 1))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// IP A uses its token
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "1.1.1.1:1111"
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// IP B should still be allowed (different limiter)
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "2.2.2.2:2222"
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
