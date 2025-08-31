package dto

import "time"

type PingDTO struct {
	Pong      bool       `json:"pong"`
	Timestamp *time.Time `json:"timestamp"`
}

type PingParams struct {
	Timestamp *time.Time `json:"timestamp"`
}
