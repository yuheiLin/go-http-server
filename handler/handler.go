package handler

import (
	"encoding/json"
	"fmt"
	"github.com/yuheiLin/go-http-server/service"
	"log"
	"net/http"
	"os"
)

type Handler interface {
	Handler1(w http.ResponseWriter, r *http.Request)
	Handler2(w http.ResponseWriter, r *http.Request)
}

type handlerImpl struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &handlerImpl{
		service: service,
	}
}

// Set up loggers
var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	// Configure loggers
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type RequestObject struct {
	R1 string `json:"r1"`
	R2 string `json:"r2"`
}
type ReturnObject struct {
	F1 string `json:"f1"`
	F2 string `json:"f2"`
}

func (h *handlerImpl) Handler1(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("handler1 called")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	p1 := r.URL.Query().Get("p1")
	fmt.Println("received p1:", p1)
	if p1 == "" {
		errorLogger.Println("missing query parameter p1")
		http.Error(w, "Bad Request: missing query parameter p1", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ReturnObject{
		F1: "f1_value_" + p1,
		F2: "f2_value",
	}); err != nil {
		log.Println("failed to encode response:", err)
	}
}

func (h *handlerImpl) Handler2(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("handle2 called")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var obj RequestObject
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		errorLogger.Println("failed to decode request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Println("received post object:", obj)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ReturnObject{
		F1: "f1_value_" + obj.R1,
		F2: "f2_value_" + obj.R2,
	}); err != nil {
		log.Println("failed to encode response:", err)
	}
}
