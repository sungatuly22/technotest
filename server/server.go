package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Server struct {
	DB *sql.DB
}

type Redirect struct {
	ActiveLink  string `json:"active_link"`
	HistoryLink string `json:"history_link"`
}

func (s Server) GetAllRedirects(w http.ResponseWriter, r *http.Request) {
	s.GetLinks()
	w.WriteHeader(http.StatusOK)
}

func (s Server) GetRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := s.GetLink(id)
	if res.ActiveLink == "" && res.HistoryLink == "" && res.Id == 0 {
		fmt.Fprintf(w, "There is no link with such id")
		return
	}
	fmt.Fprintf(w, res.ActiveLink, res.HistoryLink)
	w.WriteHeader(http.StatusOK)
}
func (s Server) DeleteRedirects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.DeleteLink(id)
	w.WriteHeader(http.StatusOK)
}
func (s Server) AddRedirects(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	infoLink := Redirects{}

	err = json.Unmarshal(data, &infoLink)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	s.AddLink(infoLink)
	w.WriteHeader(http.StatusCreated)
}

func (s Server) ChangeRedirect(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	infoLink := Redirect{}
	err = json.Unmarshal(data, &infoLink)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.UpdateLink(Redirects{Id: id, ActiveLink: infoLink.ActiveLink, HistoryLink: infoLink.HistoryLink})
	w.WriteHeader(http.StatusOK)
}
