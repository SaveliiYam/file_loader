package models

import (
	"io"
	"time"
)

type SignInRequest struct {
	Contact  string `json:"contact"`
	Password string `json:"password"`
}

type SignInResponse struct {
	APIKey string `json:"apiKey"`
}

type UploadPresignRequest struct {
	Filename string `json:"filename"`
}

type UploadPresignResponse struct {
	URLToUpload string    `json:"URLToUpload"`
	URLToLoad   string    `json:"URLToLoad"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

type ProgressReader struct {
	Reader     io.Reader
	Total      int64
	Downloaded int64
	ProgressCh chan float64
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Downloaded += int64(n)

	// Вычисляем прогресс и отправляем через канал
	progress := float64(pr.Downloaded) / float64(pr.Total) * 100
	pr.ProgressCh <- progress

	return n, err
}
