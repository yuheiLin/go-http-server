package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

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

	// routing
	r := mux.NewRouter()

	r.HandleFunc("/signup", h.CreateUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/{userID}", h.GetUserHandler).Methods(http.MethodGet)
	//r.HandleFunc("/users/{userID}", h.UpdateUserHandler).Methods(http.MethodPatch)
	//r.HandleFunc("/close", h.DeleteUserHandler).Methods(http.MethodPost)

	// serve
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server starting on port", port)
	log.Fatal(srv.ListenAndServe())
}
