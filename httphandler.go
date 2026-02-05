package envdoc

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"
)

// newDebugHandler creates the HTTP handler for GET /debug/env.
func newDebugHandler(i *Inspector) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Token auth
		if i.config.Token != "" {
			auth := r.Header.Get("Authorization")
			token := strings.TrimPrefix(auth, "Bearer ")
			if subtle.ConstantTimeCompare([]byte(token), []byte(i.config.Token)) != 1 {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Expiry check
		if !i.config.ExpiresAt.IsZero() {
			if i.clock.Now().After(i.config.ExpiresAt) {
				http.Error(w, "endpoint expired", http.StatusGone)
				return
			}
		}

		// Fresh inspection on each request
		report := i.Inspect()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	})
}
