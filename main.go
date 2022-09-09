package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// TODO(tofu345) update and delete post data.
// TODO(tofu345) Authorization

var router *mux.Router

func main() {
	router = mux.NewRouter()

	router.HandleFunc("/api/posts", getBooksView).Methods("GET")
	router.HandleFunc("/api/posts/{id}", getBookView).Methods("GET")
	router.HandleFunc("/api/posts", createPostView).Methods("POST")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe("localhost:8000", loggedRouter))
}

// func mergeMapString(mapA map[string]string, mapB map[string]string) map[string]string {
// 	for k, v := range mapB {
// 		mapA[k] = v
// 	}
// 	return mapA
// }
