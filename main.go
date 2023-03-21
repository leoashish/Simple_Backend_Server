package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type Cat struct {
	name  string
	color string
	breed string
}

type handler struct {
	db *sql.DB
}

func (h *handler) getCats(w http.ResponseWriter, r *http.Request) {
	sqlStmt := "SELECT name , color , breed from catsinfo;"
	cats := []*Cat{}
	rows, err := h.db.Query(sqlStmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		cat := &Cat{}
		err := rows.Scan(&cat.name, &cat.color, &cat.breed)
		if err != nil {
			panic(err)
		}
		cats = append(cats, cat)
	}
	json, _ := json.Marshal(cats)
	fmt.Fprint(w, json)
}

func (h *handler) AddCat(w http.ResponseWriter, r *http.Request) {
	var cat Cat
	err := json.NewDecoder(r.Body).Decode(&cat)
	if err != nil {
		return
	}
	fmt.Sprintf(cat.name, cat.color, cat.breed)
	stmt := `INSERT INTO catsinfo(name, color, breed) 
	         VALUES(?, ?, ?)`
	_, err = h.db.Exec(stmt, cat.name, cat.color, cat.breed)

	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()
	connStr := "user=postgres dbname=cats password=singoo123# host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}
	handler := handler{
		db: db,
	}
	mux.HandleFunc("/Getcats", handler.getCats)
	mux.HandleFunc("/Addcats", handler.AddCat)
	err = http.ListenAndServe("localhost:4000", mux)

	if err != nil {
		panic(err)
	}
}
