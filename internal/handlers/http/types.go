package http

import "time"

type messageResponse struct {
	ID        int32     `json:"id"`
	Message   string    `json:"message"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
