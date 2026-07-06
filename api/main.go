package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"sahil-api/cache"
	"sahil-api/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c := cache.New(5 * time.Minute)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/api/now", handlers.NowHandler(c))
	mux.HandleFunc("/api/lastfm", handlers.LastFMHandler(c))
	mux.HandleFunc("/api/anilist", handlers.AniListHandler(c))
	mux.HandleFunc("/api/letterboxd", handlers.LetterboxdHandler(c))
	mux.HandleFunc("/api/hardcover", handlers.HardcoverHandler(c))

	handler := corsMiddleware(mux)

	log.Printf("starting sahil-api on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
