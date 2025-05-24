package main

import (
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	_ "technical/docs" // Import generated docs
	"technical/posts"
)

// @title Blog API
// @version 1.0
// @description A simple blog API for managing posts
// @host localhost:8000
// @BasePath /

func main() {
	mux := http.NewServeMux()

	repo := posts.NewMapRepository()
	service := posts.NewPostService(repo)
	handler := posts.NewHandler(service)

	handler.RegisterRoutes(mux)

	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	).ServeHTTP)

	port := ":8000"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
