package structure

import "time"

type DB struct {
	Id          int
	Title       string
	Description string
	Good        int
	Bad         int
	ReplyCount  int
	View        int
	Created     time.Time
}
