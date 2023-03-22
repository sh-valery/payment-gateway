package transport

import "net/http"

func bearerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")

		switch bearerToken {
		case "Bearer first_merchant_mock_token":
			r.Header.Set("X-Merchant-Id", "merchant_one")
		case "Bearer second_merchant_mock_token":
			r.Header.Set("X-Merchant-Id", "merchant_two")
		default:
			w.WriteHeader(http.StatusUnauthorized)
		}
		next.ServeHTTP(w, r)
	})
}
