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

func newResponse(w http.ResponseWriter, responseCode int, data any, message string) {
	if responseCode == 103 {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(
		Response{
			ResponseCode: responseCode,
			Data:         data,
			Message:      message,
		})
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
	newResponse(w, 100, posts, "Post List")
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
		newResponse(w, 103, nil, "Post Not Found")
		return
	}

	newResponse(w, 100, post, "Post Detail")
}

func createPostView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post)

	postErrors, valid := post.valid()
	if valid {
		err := createPost(post)
		if err != nil {
			// TODO(tofu345) find a better way to do this
			if err.Error() == "UNIQUE constraint failed: posts.title" {
				newResponse(w, 103, nil, "Post Title Already Exists")
			} else {
				newResponse(w, 103, err, "Unexpected Error")
			}
		} else {
			newResponse(w, 100, post, "Post Created Successfully")
		}
	} else {
		newResponse(w, 103, postErrors, "Post data Invalid")
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
		newResponse(w, 103, dbError, "Error Deleting Post")
	} else {
		newResponse(w, 100, nil, "Post Deleted Successfully")
	}
}

func updatePostView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Check if post exists
	_, e := getPostById(id)
	if e != nil {
		newResponse(w, 103, nil, "Post Not Found")
	}

	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	post.ID = id // Just in case the id is passed as post params

	// Check if post is valid
	postErrors, valid := post.valid()
	if valid {
		dbError := updatePost(id, post)
		if dbError != nil {
			log.Println(err)
			newResponse(w, 103, dbError, "Error Updating Post")
		} else {
			newResponse(w, 100, post, "Post Updated Successfully")
		}
	} else {
		newResponse(w, 103, postErrors, "Post data Invalid")
	}
}
