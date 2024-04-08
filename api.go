package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(res http.ResponseWriter, status int, v any) error {
	res.WriteHeader(status)
	res.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(res).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := f(res, req); err != nil {
			WriteJSON(res, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(res http.ResponseWriter, req *http.Request) error {
	if req.Method == "POST" {
		return s.handleCreateAccount(res, req)
	}
	if req.Method == "GET" {
		return s.handleGetAccount(res, req)
	}
	if req.Method == "DELETE" {
		return s.handleDeleteAccount(res, req)
	}

	return fmt.Errorf("method not allowed: %s", req.Method)
}

func (s *APIServer) handleGetAccount(res http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)

	return WriteJSON(res, http.StatusOK, vars)
}

func (s *APIServer) handleCreateAccount(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(res http.ResponseWriter, req *http.Request) error {
	return nil
}