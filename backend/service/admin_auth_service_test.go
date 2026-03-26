package service

import (
	"testing"
	"time"

	"personal/blog/backend/model"
)

func TestAdminAuthServiceLoginAndVerify(t *testing.T) {
	authService := NewAdminAuthService("admin", "secret", "signing-key", time.Hour)

	session, token, err := authService.Login(model.AdminLoginInput{
		Username: "admin",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected login error: %v", err)
	}
	if session.Username != "admin" {
		t.Fatalf("expected admin username, got %q", session.Username)
	}
	if token == "" {
		t.Fatal("expected non-empty session token")
	}

	verified, err := authService.Verify(token)
	if err != nil {
		t.Fatalf("unexpected verify error: %v", err)
	}
	if verified.Username != "admin" {
		t.Fatalf("expected verified username admin, got %q", verified.Username)
	}
}

func TestAdminAuthServiceRejectsInvalidCredentials(t *testing.T) {
	authService := NewAdminAuthService("admin", "secret", "signing-key", time.Hour)

	_, _, err := authService.Login(model.AdminLoginInput{
		Username: "admin",
		Password: "wrong",
	})
	if err != ErrInvalidAdminCredentials {
		t.Fatalf("expected %v, got %v", ErrInvalidAdminCredentials, err)
	}
}
