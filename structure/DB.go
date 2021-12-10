package structure

import "time"

type DB struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Good        int       `json:"good"`
	Bad         int       `json:"bad"`
	ReplyCount  int       `json:"replyCount"`
	View        int       `json:"view"`
	Created     time.Time `json:"created"`
}
