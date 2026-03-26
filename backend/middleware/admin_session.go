package middleware

import (
	"context"
	"net/http"

	"personal/blog/backend/model"
	"personal/blog/backend/service"
)

type adminSessionContextKey struct{}

func RequireAdminSession(authService *service.AdminAuthService, cookieName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(cookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		session, err := authService.Verify(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), adminSessionContextKey{}, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminSessionFromContext(ctx context.Context) (model.AdminSession, bool) {
	session, ok := ctx.Value(adminSessionContextKey{}).(model.AdminSession)
	return session, ok
}
