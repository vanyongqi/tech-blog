package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"personal/blog/backend/model"
)

var ErrInvalidAdminCredentials = errors.New("invalid admin credentials")
var ErrInvalidAdminSession = errors.New("invalid admin session")

type AdminAuthService struct {
	username   string
	password   string
	secret     []byte
	sessionTTL time.Duration
}

func NewAdminAuthService(username, password, secret string, sessionTTL time.Duration) *AdminAuthService {
	return &AdminAuthService{
		username:   strings.TrimSpace(username),
		password:   password,
		secret:     []byte(secret),
		sessionTTL: sessionTTL,
	}
}

func (s *AdminAuthService) Login(input model.AdminLoginInput) (model.AdminSession, string, error) {
	if strings.TrimSpace(input.Username) != s.username || input.Password != s.password {
		return model.AdminSession{}, "", ErrInvalidAdminCredentials
	}

	session := model.AdminSession{
		Username:  s.username,
		ExpiresAt: time.Now().Add(s.sessionTTL).UTC(),
	}
	return session, s.signSession(session), nil
}

func (s *AdminAuthService) Verify(token string) (model.AdminSession, error) {
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(token))
	if err != nil {
		return model.AdminSession{}, ErrInvalidAdminSession
	}

	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 {
		return model.AdminSession{}, ErrInvalidAdminSession
	}

	payload := parts[0] + "|" + parts[1]
	expectedSig := s.sign(payload)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return model.AdminSession{}, ErrInvalidAdminSession
	}

	expiresAt, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return model.AdminSession{}, ErrInvalidAdminSession
	}
	if time.Now().UTC().After(expiresAt) {
		return model.AdminSession{}, ErrInvalidAdminSession
	}

	return model.AdminSession{
		Username:  parts[0],
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AdminAuthService) signSession(session model.AdminSession) string {
	payload := fmt.Sprintf("%s|%s", session.Username, session.ExpiresAt.Format(time.RFC3339))
	token := payload + "|" + s.sign(payload)
	return base64.RawURLEncoding.EncodeToString([]byte(token))
}

func (s *AdminAuthService) sign(payload string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
