package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
)

type VisitorIdentity struct {
	ID          string
	DisplayName string
	IP          string
}

type visitorContextKey struct{}

func CaptureVisitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity := buildVisitorIdentity(r)
		ctx := context.WithValue(r.Context(), visitorContextKey{}, identity)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MustVisitorIdentity(ctx context.Context) VisitorIdentity {
	identity, ok := VisitorIdentityFromContext(ctx)
	if !ok {
		return VisitorIdentity{}
	}
	return identity
}

func VisitorIdentityFromContext(ctx context.Context) (VisitorIdentity, bool) {
	identity, ok := ctx.Value(visitorContextKey{}).(VisitorIdentity)
	return identity, ok
}

func buildVisitorIdentity(r *http.Request) VisitorIdentity {
	ip := extractClientIP(r)
	userAgent := strings.TrimSpace(r.UserAgent())
	raw := ip + "\n" + userAgent
	digest := sha256.Sum256([]byte(raw))
	id := hex.EncodeToString(digest[:])
	label := "访客-" + strings.ToUpper(id[:8])
	return VisitorIdentity{
		ID:          id,
		DisplayName: label,
		IP:          ip,
	}
}

func extractClientIP(r *http.Request) string {
	if forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
