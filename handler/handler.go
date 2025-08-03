package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/yuheiLin/go-http-server/customerror"
	"github.com/yuheiLin/go-http-server/model"
	"github.com/yuheiLin/go-http-server/service"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Handler interface {
	GetUserHandler(w http.ResponseWriter, r *http.Request)
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
	//DeleteUserHandler(w http.ResponseWriter, r *http.Request)
}

type handlerImpl struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &handlerImpl{
		service: service,
	}
}

type APIResponse struct {
	Message string      `json:"message,omitempty"`
	User    *model.User `json:"user,omitempty"`
	Cause   string      `json:"cause,omitempty"`
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

func validateLettersAndDigits(input string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", input)
	return match
}

func (h *handlerImpl) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("GetUserHandler called")

	authHeader := r.Header.Get("Basic")
	authDecoded, _ := base64.StdEncoding.DecodeString(authHeader)
	authDecodedString := string(authDecoded)
	idpass := strings.Split(authDecodedString, ":")
	if len(idpass) != 2 {
		errorLogger.Println("Invalid Basic Auth header format")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.service.VerifyUser(idpass[0], idpass[1]); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := APIResponse{
			Message: "Authentication failed",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("failed to encode response:", err)
		}
		return
	}

	// get user ID from path parameters
	vars := mux.Vars(r)
	userID := vars["userID"]
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := APIResponse{
			Message: "Get User failed",
			Cause:   "Required user_id",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("failed to encode response:", err)
		}
		return
	}

	u, err := h.service.GetUser(userID)
	if err != nil {
		if errors.Is(err, customerror.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			response := APIResponse{
				Message: "No user found",
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Println("failed to encode response:", err)
			}
			return
		}
		errorLogger.Println("failed to get user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	u.FillNickname()
	response := APIResponse{
		Message: "User details by user_id",
		User:    u,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to encode response:", err)
	}
}

func (h *handlerImpl) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	infoLogger.Println("CreateUserHandler called")
	w.Header().Set("Content-Type", "application/json")

	var userRequest model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		errorLogger.Println("failed to decode request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// validate user request
	if userRequest.UserID == "" || userRequest.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := APIResponse{
			Message: "Account creation failed",
			Cause:   "Required user_id and password",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("failed to encode response:", err)
		}
		return
	}
	if len(userRequest.Password) < 8 || len(userRequest.Password) > 20 || len(userRequest.UserID) < 6 || len(userRequest.UserID) > 20 {
		w.WriteHeader(http.StatusBadRequest)
		response := APIResponse{
			Message: "Account creation failed",
			Cause:   "Input length is incorrect",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("failed to encode response:", err)
		}
		return
	}
	// TODO: password allows special characters
	if !validateLettersAndDigits(userRequest.UserID) || !validateLettersAndDigits(userRequest.Password) {
		w.WriteHeader(http.StatusBadRequest)
		response := APIResponse{
			Message: "Account creation failed",
			Cause:   "Incorrect character pattern",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("failed to encode response:", err)
		}
		return
	}

	createdUser, err := h.service.CreateUser(userRequest.UserID, userRequest.Password)
	if err != nil {
		if errors.Is(err, customerror.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			response := APIResponse{
				Message: "Account creation failed",
				Cause:   "Already same user_id is used",
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Println("failed to encode response:", err)
			}
			return
		}
		errorLogger.Println("failed to create user:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	createdUser.FillNickname()
	response := APIResponse{
		Message: "Account successfully created",
		User:    createdUser,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to encode response:", err)
	}
}

//func (h *handlerImpl) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
//	infoLogger.Println("DeleteUserHandler called")
//
//	if r.Method != http.MethodDelete {
//		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//		return
//	}
//
//	// get user ID from path parameters
//	vars := mux.Vars(r)
//	userID := vars["userID"]
//	if userID == "" {
//		errorLogger.Println("missing path parameter userID")
//		http.Error(w, "Bad Request: missing path parameter userID", http.StatusBadRequest)
//		return
//	}
//
//	if err := h.service.DeleteUser(userID); err != nil {
//		errorLogger.Println("failed to delete user:", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}
