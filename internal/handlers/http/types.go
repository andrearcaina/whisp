package http

import (
	"time"
)

type messageResponse struct {
	ID        int32     `json:"id"`
	Message   *string   `json:"message,omitempty"`
	ImageUrl  *string   `json:"image_url,omitempty"`
	GifUrl    *string   `json:"gif_url,omitempty"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type klipyGIFResponse struct {
	Results []struct {
		ID           string `json:"id"`
		Title        string `json:"title"`
		MediaFormats struct {
			GIF struct {
				URL string `json:"url"`
			} `json:"gif"`
		} `json:"media_formats"`
	} `json:"results"`
}

type klipyGIFTrimmedResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	GifUrl string `json:"gif_url"`
}
