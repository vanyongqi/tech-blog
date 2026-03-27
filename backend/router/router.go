package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"personal/blog/backend/controller"
	"personal/blog/backend/middleware"
	"personal/blog/backend/service"
)

func NewHandler(blogController *controller.BlogController, adminController *controller.AdminController, authService *service.AdminAuthService, adminCookieName string, staticDir string) http.Handler {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/home", methodGuard(http.MethodGet, blogController.GetHome))
	apiMux.HandleFunc("/api/posts", methodGuard(http.MethodGet, blogController.ListPosts))
	apiMux.HandleFunc("/api/posts/", blogController.HandlePostRoute)
	apiMux.HandleFunc("/api/assets/", methodGuard(http.MethodGet, blogController.GetAsset))

	adminPublicMux := http.NewServeMux()
	adminPublicMux.HandleFunc("/api/admin/login", methodGuard(http.MethodPost, adminController.Login))

	adminProtectedMux := http.NewServeMux()
	adminProtectedMux.HandleFunc("/api/admin/logout", methodGuard(http.MethodPost, adminController.Logout))
	adminProtectedMux.HandleFunc("/api/admin/session", methodGuard(http.MethodGet, adminController.GetSession))
	adminProtectedMux.HandleFunc("/api/admin/assets", methodGuard(http.MethodPost, adminController.UploadAsset))
	adminProtectedMux.HandleFunc("/api/admin/posts", adminController.HandlePostsRoute)
	adminProtectedMux.HandleFunc("/api/admin/posts/", adminController.HandlePostsRoute)
	adminProtectedMux.HandleFunc("/api/admin/projects", adminController.HandleProjectsRoute)
	adminProtectedMux.HandleFunc("/api/admin/projects/", adminController.HandleProjectsRoute)
	adminProtectedMux.HandleFunc("/api/admin/videos", adminController.HandleVideosRoute)
	adminProtectedMux.HandleFunc("/api/admin/videos/", adminController.HandleVideosRoute)

	rootAPIMux := http.NewServeMux()
	rootAPIMux.Handle("/api/admin/login", withCORS(adminPublicMux))
	rootAPIMux.Handle("/api/admin/logout", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/admin/session", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/admin/assets", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/admin/posts", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/admin/posts/", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/admin/projects", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/admin/projects/", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/admin/videos", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/admin/videos/", withCORS(middleware.RequireAdminSession(authService, adminCookieName, adminProtectedMux)))
	rootAPIMux.Handle("/api/home", middleware.CaptureVisitor(withCORS(apiMux)))
	rootAPIMux.Handle("/api/posts", middleware.CaptureVisitor(withCORS(apiMux)))
	rootAPIMux.Handle("/api/posts/", middleware.CaptureVisitor(withCORS(apiMux)))
	rootAPIMux.Handle("/api/assets/", withCORS(apiMux))

	if strings.TrimSpace(staticDir) == "" {
		rootMux := http.NewServeMux()
		rootMux.Handle("/api/", rootAPIMux)
		rootMux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("personal blog api is running"))
		})
		return rootMux
	}

	return &spaHandler{
		apiHandler: rootAPIMux,
		staticDir:  staticDir,
	}
}

type spaHandler struct {
	apiHandler http.Handler
	staticDir  string
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		h.apiHandler.ServeHTTP(w, r)
		return
	}

	cleanPath := strings.TrimPrefix(filepath.Clean("/"+r.URL.Path), "/")
	target := filepath.Join(h.staticDir, cleanPath)
	if info, err := os.Stat(target); err == nil && !info.IsDir() {
		http.ServeFile(w, r, target)
		return
	}

	http.ServeFile(w, r, filepath.Join(h.staticDir, "index.html"))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Vary", "Origin")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func methodGuard(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
