package http

import (
	"time"
)

type messageResponse struct {
	ID        int32     `json:"id"`
	Message   string    `json:"message"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type tenorGIFResponse struct {
	Results []struct {
		ID                 string `json:"id"`
		ContentDescription string `json:"content_description"`
		MediaFormats       struct {
			GIF struct {
				URL string `json:"url"`
			} `json:"gif"`
		} `json:"media_formats"`
	} `json:"results"`
}

type tenorGIFTrimmedResponse struct {
	ID     string `json:"id"`
	Desc   string `json:"desc"`
	GifUrl string `json:"gif_url"`
}
