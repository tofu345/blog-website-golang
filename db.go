package main

import (
	"database/sql"
	"log"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		log.Fatal("Error Connecting to DB: ", err)
	}
	createTable()
}

func createTable() {
	post_table := `CREATE TABLE posts (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "title" TEXT UNIQUE,
        "content" TEXT,
        "author" TEXT,
        "views" INT);`

	query, err := db.Prepare(post_table)
	if err != nil {
		log.Println(err)
	} else {
		query.Exec()
	}
}

func createPost(post Post) error {
	records := `INSERT INTO posts(title, content, author, views) VALUES (?, ?, ?, ?)`
	query, err := db.Prepare(records)
	if err != nil {
		log.Fatal(err)
	}

	_, err = query.Exec(post.Title, post.Content, post.Author, 0)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func getPostById(object_id int) (Post, error) {
	query := db.QueryRow("SELECT * FROM posts WHERE id = ?", object_id)
	var id int
	var title string
	var content string
	var author string
	var views int
	err := query.Scan(&id, &title, &content, &author, &views)
	if err != nil {
		return Post{}, err
	}
	post := newPost(id, title, content, author, views)
	return post, nil
}

func deletePost(object_id int) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = ?", object_id)

	if err != nil {
		return err
	} else {
		return nil
	}
}

// func updatePost(object_id int, newPostData Post) error {

// }
