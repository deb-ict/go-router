package authentication

import "net/http"

type authenticationMiddleware struct {
}

func (m *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: Handle authentication
		next.ServeHTTP(w, r)
	})
}
