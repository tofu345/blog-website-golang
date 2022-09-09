package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct {
	ResponseCode int    `json:"responseCode"`
	Message      string `json:"message"`
	Data         any    `json:"data"`
}

func newResponse(code int, data any, message string) Response {
	return Response{
		ResponseCode: code,
		Data:         data,
		Message:      message,
	}
}

func getBooksView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	record, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	defer record.Close()

	var posts []Post
	for record.Next() {
		var id int
		var title string
		var content string
		var author string
		var views int
		record.Scan(&id, &title, &content, &author, &views)
		posts = append(posts, newPost(id, title, content, author, views))
	}
	json.NewEncoder(w).Encode(newResponse(100, posts, "Post List"))
}

func getBookView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	post, err := getPostById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newResponse(103, nil, "Post Not Found"))
		return
	}

	json.NewEncoder(w).Encode(newResponse(100, post, "Post Detail"))
}

func createPostView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	post_errors, valid := post.valid()
	if valid {
		err := createPost(post)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: posts.title" {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(newResponse(103, nil, "Post Title Already Exists"))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(newResponse(103, err, "Unexpected Error"))
			}
		} else {
			json.NewEncoder(w).Encode(newResponse(100, post, "Post Created Successfully"))
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(newResponse(103, post_errors, "Post data Invalid"))
	}
}

func deletePostView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	dbError := deletePost(id)
	if dbError != nil {
		json.NewEncoder(w).Encode(newResponse(103, dbError, "Error Deleting Post"))
	} else {
		json.NewEncoder(w).Encode(newResponse(100, nil, "Post Deleted Successfully"))
	}
}
