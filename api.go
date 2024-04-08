package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	storage Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(res http.ResponseWriter, req *http.Request) error {
	if req.Method == "POST" {
		return s.handleCreateAccount(res, req)
	}
	if req.Method == "GET" {
		return s.handleGetAccounts(res, req)
	}
	if req.Method == "DELETE" {
		return s.handleDeleteAccount(res, req)
	}

	return fmt.Errorf("method not allowed: %s", req.Method)
}

func (s *APIServer) handleGetAccountByID(res http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)

	return WriteJSON(res, http.StatusOK, vars)
}

func (s *APIServer) handleGetAccounts(res http.ResponseWriter, req *http.Request) error {
	accounts, err := s.storage.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(res, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(res http.ResponseWriter, req *http.Request) error {
	requestArgs := new(CreateAccountRequest)
	if err := json.NewDecoder(req.Body).Decode(requestArgs); err != nil {
		return err
	}

	account := NewAccount(requestArgs.FirstName, requestArgs.LastName)
	if err := s.storage.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(res, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(res http.ResponseWriter, req *http.Request) error {
	return nil
}

func WriteJSON(res http.ResponseWriter, status int, v any) error {
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)
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