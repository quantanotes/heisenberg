package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	db *DB
}

func NewServer(db *DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) StartServer() {
	http.HandleFunc("/newcollection/", s.handleNewCollection)
	http.HandleFunc("/delcollection/", s.handleDelCollection)
	http.HandleFunc("/put/", s.handlePut)
	http.HandleFunc("/get/", s.handleGet)
	http.HandleFunc("/del/", s.handleDel)
	fmt.Println("Starting server on :420")
	http.ListenAndServe(":420", nil)
}

func (s *Server) handleNewCollection(w http.ResponseWriter, r *http.Request) {
	var b *struct {
		Collection string `json:"collection"`
	}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = s.db.NewCollection(b.Collection)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

func (s *Server) handleDelCollection(w http.ResponseWriter, r *http.Request) {
	var b *struct {
		Collection string `json:"collection"`
	}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = s.db.NewCollection(b.Collection)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	b := &struct {
		Key        string    `json:"key"`
		Vec        []float32 `json:"vec"`
		Meta       any       `json:"meta"`
		Collection string    `json:"collection"`
	}{}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	pair := pair{b.Key, value{b.Vec, b.Meta}}

	err = s.db.Put(pair, b.Collection)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	b := &struct {
		Key        string `json:"key"`
		Collection string `json:"collection"`
	}{}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	pair, err := s.db.Get(b.Key, b.Collection)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data, err := json.Marshal(pair)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(data)
}

func (s *Server) handleDel(w http.ResponseWriter, r *http.Request) {
	b := &struct {
		Key        string `json:"key"`
		Collection string `json:"collection"`
	}{}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = s.db.Delete(b.Key, b.Collection)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}
