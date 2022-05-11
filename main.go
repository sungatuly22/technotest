package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"technotest/server"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./redirects.db")
	if err != nil {
		log.Fatalf(err.Error())
	}
	s := server.NewDb(db)

	file, err := ioutil.ReadFile("links.json")
	if err != nil {
		log.Fatalf(err.Error())
	}

	data := []server.Redirect{}

	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatalf(err.Error())
	}
	for i := 0; i < len(data); i++ {
		s.AddLink(server.Redirects{ActiveLink: data[i].ActiveLink, HistoryLink: data[i].HistoryLink})
	}
	r := mux.NewRouter()
	r.HandleFunc("/admin/redirects", s.GetAllRedirects).Methods("GET")
	r.HandleFunc("/admin/redirects/{id}", s.GetRedirect).Methods("GET")
	r.HandleFunc("/admin/redirects/{id}", s.DeleteRedirects).Methods("DELETE")
	r.HandleFunc("/admin/redirects", s.AddRedirects).Methods("POST")
	r.HandleFunc("/admin/redirects/{id}", s.ChangeRedirect).Methods("PATCH")

	if err := http.ListenAndServe("localhost:8080", r); err != nil {
		log.Fatalf(err.Error())
	}
}
