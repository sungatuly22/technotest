package server

import (
	"database/sql"
	"fmt"
	"log"
)

type Redirects struct {
	Id          int    `json:"Id"`
	ActiveLink  string `json:"active_link"`
	HistoryLink string `json:"history_link"`
}

func (s *Server) GetLinks() {
	var (
		id          int
		activeLink  string
		historyLink string
	)
	links := []Redirects{}
	rows, err := s.DB.Query(`
			SELECT * FROM redirects
			`)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&id, &activeLink, &historyLink)
		link := Redirects{Id: id, ActiveLink: activeLink, HistoryLink: historyLink}
		links = append(links, link)
	}
	fmt.Println(links)
}

func (s *Server) GetLink(id int) Redirects {
	link := Redirects{}
	row := s.DB.QueryRow("SELECT id, activelink, historylink FROM redirects WHERE id=?", id)
	row.Scan(&link.Id, &link.ActiveLink, &link.HistoryLink)
	return link
}

func (s *Server) AddLink(link Redirects) {
	stmt, err := s.DB.Prepare(`
				INSERT INTO redirects (activelink, historylink) values (?, ?)
				`)
	if err != nil {
		log.Fatalf(err.Error())
	}
	stmt.Exec(link.ActiveLink, link.HistoryLink)
}

func NewDb(db *sql.DB) *Server {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS redirects (
								id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
								activelink TEXT, 
								historylink TEXT);
							`)
	if err != nil {
		log.Fatalf(err.Error())
	}
	stmt.Exec()
	return &Server{DB: db}
}

func (s *Server) DeleteLink(id int) {
	_, err := s.DB.Exec("DELETE FROM redirects WHERE id=$1", id)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (s *Server) UpdateLink(link Redirects) {
	_, err := s.DB.Exec("UPDATE redirects set activelink = $1, historylink = $2 WHERE id = $3", link.ActiveLink, link.HistoryLink, link.Id)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
