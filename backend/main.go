package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
	"personal/blog/backend/controller"
	"personal/blog/backend/dao"
	"personal/blog/backend/router"
	"personal/blog/backend/service"
)

func main() {
	ctx := context.Background()
	dbPath := envOrDefault("BLOG_DB_PATH", filepath.Join("storage", "blog.db"))
	db, err := dao.InitSQLite(ctx, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repository := dao.NewSQLiteRepository(db)
	blogService := service.NewBlogService(repository)
	blogController := controller.NewBlogController(blogService)
	adminUsername := envOrDefault("BLOG_ADMIN_USER", "admin")
	adminPassword := envRequired("BLOG_ADMIN_PASSWORD")
	adminSecret := envRequired("BLOG_ADMIN_SECRET")
	adminCookieName := envOrDefault("BLOG_ADMIN_COOKIE_NAME", "blog_admin_session")
	adminCookieSecure := envBool("BLOG_ADMIN_COOKIE_SECURE", false)
	adminSessionTTL := envDurationHours("BLOG_ADMIN_SESSION_HOURS", 72)

	adminAuthService := service.NewAdminAuthService(adminUsername, adminPassword, adminSecret, adminSessionTTL)
	adminContentService := service.NewAdminContentService(repository)
	adminProjectService := service.NewAdminProjectService(repository)
	adminVideoService := service.NewAdminVideoService(repository)
	adminController := controller.NewAdminController(adminAuthService, adminContentService, adminProjectService, adminVideoService, adminCookieName, adminCookieSecure)

	addr := envOrDefault("BLOG_ADDR", ":8080")
	staticDir := stringsOrEmpty("FRONTEND_DIST")

	handler := router.NewHandler(blogController, adminController, adminAuthService, adminCookieName, staticDir)
	log.Printf("blog backend listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func stringsOrEmpty(key string) string {
	return os.Getenv(key)
}

func envRequired(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}

func envBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	switch strings.ToLower(value) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func envDurationHours(key string, fallbackHours int) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return time.Duration(fallbackHours) * time.Hour
	}
	hours, err := strconv.Atoi(value)
	if err != nil || hours <= 0 {
		return time.Duration(fallbackHours) * time.Hour
	}
	return time.Duration(hours) * time.Hour
}
