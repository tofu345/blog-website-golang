package main

import (
	"strconv"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Views   int    `json:"views"`
}

func (post *Post) valid() (map[string]string, bool) {
	errors := map[string]string{}

	if post.Title == "" {
		errors["title"] = "This field is required"
	}

	if post.Content == "" {
		errors["content"] = "This field is required"
	}

	if post.Author == "" {
		errors["author"] = "This field is required"
	}

	if len(errors) > 0 {
		return errors, false
	} else {
		return errors, true
	}
}

func (post *Post) format() map[string]string {
	return map[string]string{
		"id":      strconv.Itoa(post.ID),
		"title":   post.Title,
		"content": post.Content,
		"author":  post.Author,
		"views":   strconv.Itoa(post.Views),
	}
}

// func (post *Post) print() {
// 	output := fmt.Sprintf("<Post %v: %v by %v>", post.ID, post.Title, post.Author)
// 	fmt.Println(output)
// }

func newPost(id int, title string, content string, author string, views int) Post {
	return Post{
		ID:      id,
		Title:   title,
		Content: content,
		Author:  author,
		Views:   views,
	}
}

func FormatPost(posts []Post) []map[string]string {
	output := []map[string]string{}
	for _, post := range posts {
		output = append(output, post.format())
	}
	return output
}
