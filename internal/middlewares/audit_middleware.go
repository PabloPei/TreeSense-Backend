package middlewares

import (
	"net/http"
	"strings"
	"context"
	"regexp"

)

func NewAuditMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := context.WithValue(r.Context(), "audit", true)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func inferActionName(method, path string) string {
	methodMap := map[string]string{
		"POST":   "create",
		"PUT":    "update",
		"PATCH":  "update",
		"DELETE": "delete",
		"GET":    "read",
	}

	cleanPath := strings.ToLower(path)
	re := regexp.MustCompile(`^/api/v\d+`)
	cleanPath = re.ReplaceAllString(cleanPath, "")
	cleanPath = strings.ReplaceAll(cleanPath, "/", "_")
	cleanPath = strings.ReplaceAll(cleanPath, "{", "")
	cleanPath = strings.ReplaceAll(cleanPath, "}", "")

	action := methodMap[method]
	if action == "" {
		action = "unknown"
	}

	return action + cleanPath
}
