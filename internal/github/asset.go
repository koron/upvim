package github

import "time"

type Asset struct {
	Name        string
	State       string
	Size        uint64
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DownloadURL string    `json:"browser_download_url"`
}
