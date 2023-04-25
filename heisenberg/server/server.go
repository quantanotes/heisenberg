package server

import (
	"encoding/json"
	"fmt"
	"heisenberg/storage"
	"net/http"
)

type Server struct {
	db *storage.DB
}

func NewServer(db *storage.DB) *Server {
	return &Server{
		db,
	}
}

func (s *Server) Run() {
	http.HandleFunc("/newcollection/", s.handleNewCollection)
	http.HandleFunc("/delcollection/", s.handleDelCollection)
	http.HandleFunc("/put/", s.handlePut)
	http.HandleFunc("/get/", s.handleGet)
	http.HandleFunc("/del/", s.handleDel)
	http.HandleFunc("/search/", s.handleSearch)
	fmt.Println("Starting server on :420")
	http.ListenAndServe("localhost:420", nil)
}

func (s *Server) handleNewCollection(w http.ResponseWriter, r *http.Request) {
	b := &struct {
		Name  string `json:"name"`
		Dim   int    `json:"dim"`
		Size  int    `json:"size"`
		Space string `json:"space"`
		M     int    `json:"m"`
		Ef    int    `json:"ef"`
	}{}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = s.db.NewCollection(b.Name, b.Dim, b.Size, b.Space, b.M, b.Ef)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

func (s *Server) handleDelCollection(w http.ResponseWriter, r *http.Request) {
	b := &struct {
		Name string `json:"name"`
	}{}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = s.db.DeleteCollection(b.Name)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
}

func (s *Server) handlePut(w http.ResponseWriter, r *http.Request) {
	b := &struct {
		Key        string      `json:"key"`
		Vec        []float32   `json:"vec"`
		Meta       interface{} `json:"meta"`
		Collection string      `json:"collection"`
	}{}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = s.db.Put(b.Key, b.Vec, b.Meta, b.Collection)
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

	value, err := s.db.Get(b.Key, b.Collection)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	pair := storage.Pair{
		b.Key,
		*value,
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

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	b := &struct {
		Vec        []float32 `json:"vec"`
		K          int       `json:"k"`
		Collection string    `json:"collection"`
	}{}

	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res, err := s.db.Search(b.Vec, b.K, b.Collection)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(data)
}
