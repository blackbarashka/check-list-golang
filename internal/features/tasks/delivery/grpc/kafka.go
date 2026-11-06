package grpc

import "time"

type Event struct {
	Action string    `json:"action"`
	Time   time.Time `json:"time"`
}
