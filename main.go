package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yuheiLin/go-http-server/apiclient"
	"github.com/yuheiLin/go-http-server/handler"
	"github.com/yuheiLin/go-http-server/repository"
	"github.com/yuheiLin/go-http-server/service"
)

func main() {
	// get envs
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // Default to 10000 if PORT env var is not set
	}

	client := apiclient.New()
	repo := repository.New()
	svc := service.New(repo, client)
	h := handler.New(svc)
	http.HandleFunc("/h1", h.Handler1)
	http.HandleFunc("/h2", h.Handler2)
	log.Println("Server starting on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
