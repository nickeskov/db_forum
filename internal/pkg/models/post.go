package models

import "time"

//easyjson:json
type Post struct {
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	Created  time.Time `json:"created,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	ID       int64     `json:"id,omitempty"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Parent   int64     `json:"parent,omitempty"`
	Thread   int32     `json:"thread,omitempty"`
}

//easyjson:json
type Posts []Post

//easyjson:json
type PostUpdate struct {
	Description string `json:"description,omitempty"`
	Message     string `json:"message,omitempty"`
}

//easyjson:json
type PostUpdates []PostUpdate

//easyjson:json
type PostFullInfo struct {
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Post   *Post   `json:"post,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}

//easyjson:json
type PostsFulls []PostFullInfo

func (post Post) Validate() error {
	if post.Message == "" {
		return ErrValidation
	}
	return nil
}
