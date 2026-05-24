package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"pgreader/handlers"
	"pgreader/services"
)

func main() {
	dataFile := getenv("PG_READER_DATA_FILE", filepath.Join("..", "data", "articles.json"))
	addr := getenv("PG_READER_ADDR", ":8080")

	store, err := services.NewStore(dataFile)
	if err != nil {
		log.Fatalf("store init: %v", err)
	}

	articleSvc := services.NewArticleService(store, services.NewScraper())

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()
	if err := articleSvc.EnsureSynced(ctx); err != nil {
		log.Printf("initial sync skipped: %v", err)
	}

	mux := http.NewServeMux()
	handlers.NewHTTPHandler(articleSvc).Register(mux)

	// Frontend dev server can call backend directly during local development.
	muxWithCORS := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})

	log.Printf("pg-reader backend listening on %s", addr)
	if err := http.ListenAndServe(addr, muxWithCORS); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
