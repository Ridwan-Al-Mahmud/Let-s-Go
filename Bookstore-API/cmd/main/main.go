package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Ridwan-Al-Mahmud/Go-Bookstore/pkg/routes"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9010", r))
}
