package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yuheiLin/go-http-server/model"
	"github.com/yuheiLin/go-http-server/service"
	"log"
	"net/http"
	"os"
)

type Handler interface {
	GetHandler(w http.ResponseWriter, r *http.Request)
	PostHandler(w http.ResponseWriter, r *http.Request)
	GetUserHandler(w http.ResponseWriter, r *http.Request)
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
	DeleteUserHandler(w http.ResponseWriter, r *http.Request)
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

func (h *handlerImpl) GetHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("GetHandler called")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// get query parameters
	p1 := r.URL.Query().Get("p1")
	fmt.Println("received p1:", p1)
	if p1 == "" {
		errorLogger.Println("missing query parameter p1")
		http.Error(w, "Bad Request: missing query parameter p1", http.StatusBadRequest)
		return
	}

	// get path parameters
	vars := mux.Vars(r)
	id1 := vars["id1"]
	id2 := vars["id2"]

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ReturnObject{
		F1: fmt.Sprintf("f1_value_id1_%s_id2_%s_p1_%s", id1, id2, p1),
		F2: "f2_value",
	}); err != nil {
		log.Println("failed to encode response:", err)
	}
}

func (h *handlerImpl) PostHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("PostHandler called")

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

func (h *handlerImpl) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("GetUserHandler called")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// get user ID from path parameters
	vars := mux.Vars(r)
	userID := vars["userID"]
	if userID == "" {
		errorLogger.Println("missing path parameter userID")
		http.Error(w, "Bad Request: missing path parameter userID", http.StatusBadRequest)
		return
	}

	u, err := h.service.GetUser(userID)
	if err != nil {
		errorLogger.Println("failed to get user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Println("failed to encode response:", err)
	}
}

func (h *handlerImpl) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("CreateUserHandler called")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errorLogger.Println("failed to decode request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(user.ID, user.Name); err != nil {
		errorLogger.Println("failed to create user:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handlerImpl) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("DeleteUserHandler called")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// get user ID from path parameters
	vars := mux.Vars(r)
	userID := vars["userID"]
	if userID == "" {
		errorLogger.Println("missing path parameter userID")
		http.Error(w, "Bad Request: missing path parameter userID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(userID); err != nil {
		errorLogger.Println("failed to delete user:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
