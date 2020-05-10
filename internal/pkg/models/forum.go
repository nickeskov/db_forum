package models

//easyjson:json
type Forum struct {
	ID      int64  `json:"-"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	User    string `json:"user"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int64  `json:"threads,omitempty"`
}

//easyjson:json
type Forums []Forum

func (forum Forum) Validate() error {
	return validateSlug(forum.Slug)
}
