package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type Cat struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Breed string `json:"breed"`
}

type handler struct {
	db *sql.DB
}

func (h *handler) getCats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	sqlStmt := "SELECT name , color , breed from catsInfo;"
	cats := []*Cat{}
	rows, err := h.db.Query(sqlStmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		cat := &Cat{}
		err := rows.Scan(&cat.Name, &cat.Color, &cat.Breed)
		if err != nil {
			panic(err)
		}
		cats = append(cats, cat)
	}
	if err != nil {
		panic(err)
	}
	json, err := json.Marshal(&cats)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
	w.Write(json)
}

func (h *handler) AddCat(w http.ResponseWriter, r *http.Request) {
	var cat Cat
	err := json.NewDecoder(r.Body).Decode(&cat)
	if err != nil {
		return
	}

	stmt := `INSERT INTO catsInfo(name, color, breed)
	         VALUES ($1,$2,$3);`
	_, err = h.db.Exec(stmt, cat.Name, cat.Color, cat.Breed)

	if err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()
	connStr := "user=postgres dbname=cats password=notMyRealPassword host=localhost sslmode=disable"
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
