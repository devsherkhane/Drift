package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adds industry-standard security headers to every HTTP response.
// These headers protect against common web vulnerabilities like clickjacking,
// MIME-type sniffing, and cross-site scripting.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME-type sniffing — browser must respect Content-Type
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent the page from being embedded in iframes (clickjacking protection)
		c.Header("X-Frame-Options", "DENY")

		// Enable browser XSS filter and block rendering if attack detected
		c.Header("X-XSS-Protection", "1; mode=block")

		// Control how much referrer info is sent with requests
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Restrict access to browser features like camera, microphone, etc.
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		c.Next()
	}
}
