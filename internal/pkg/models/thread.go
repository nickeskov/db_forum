package models

import "time"

//easyjson:json
type Thread struct {
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Message string    `json:"message"`
	ID      int32     `json:"id,omitempty"`
	Forum   string    `json:"forum,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Votes   int32     `json:"votes,omitempty"`
}

//easyjson:json
type Threads []Thread

func (thread Thread) Validate() error {
	if thread.Slug != "" {
		return validateSlug(thread.Slug)
	}
	return nil
}
