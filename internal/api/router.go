package api

import (
	"log"
	"net/http"
)

// Router menyiapkan router untuk API HTTP
func (h *Handler) SetupRouter() http.Handler {
	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("/api/chat", h.HandleChat)
	mux.HandleFunc("/api/documents", h.HandleAddDocument)
	mux.HandleFunc("/api/conversations", h.HandleGetConversation)

	// Middleware untuk logging dan CORS
	return logMiddleware(corsMiddleware(mux))
}

// logMiddleware mencatat permintaan HTTP
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware menambahkan header CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
