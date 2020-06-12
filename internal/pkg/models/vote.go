package models

//easyjson:json
type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int16  `json:"voice"`
}

//easyjson:json
type Votes []Vote
