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
	Time        time.Time `json:"time"`
	Created     string    `json:"created"`
	OwnerId     int       `json:"ownerId"`
}
