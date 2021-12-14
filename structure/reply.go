package structure

import "time"

type Reply struct {
	Id          int       `json:"id"`
	OwnerId     int       `json:"ownerId"`
	Description string    `json:"description"`
	PostId      int       `json:"postId"`
	Time        time.Time `json:"time"`
	Created     string    `json:"created"`
}
