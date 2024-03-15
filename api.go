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
	store Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountById))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func(s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	} else if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	} else if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

// GET /acount
func(s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func(s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	// account := NewAccount("Amirabbas", "Ghasemi")
	return WriteJSON(w, http.StatusOK, &Account{

	});
}

func(s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
    decoder := json.NewDecoder(r.Body)
    t := &CreateAccountRequest{}
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }
	account := NewAccount(t.FirstName, t.LastName)

	lastInsertId, err := s.store.CreateAccount(account)

	if err != nil {
		return err
	}

	account.ID = lastInsertId

	return WriteJSON(w, http.StatusCreated, account)
}

func(s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil;
}

func(s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil;
}


func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func (http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
