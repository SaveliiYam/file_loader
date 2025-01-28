package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"ss/internal/config"
	"ss/internal/models"
	"time"
)

func SignInImmediate() (string, error) {
	data := models.SignInRequest{
		Contact:  config.Contact,
		Password: config.Password,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		config.BasePath+"auth/v1/immediate_sign_in",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var signInResponse models.SignInResponse
	if err := json.NewDecoder(resp.Body).Decode(&signInResponse); err != nil {
		return "", err
	}

	return signInResponse.APIKey, nil
}

func UploadPresign(file, apiKey string) (models.UploadPresignResponse, error) {
	data := models.UploadPresignRequest{
		Filename: file,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return models.UploadPresignResponse{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		config.BasePath+"storage/v1/upload_presign",
		bytes.NewReader(body),
	)
	if err != nil {
		return models.UploadPresignResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return models.UploadPresignResponse{}, err
	}
	defer resp.Body.Close()

	var uploadPresignResponse models.UploadPresignResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadPresignResponse); err != nil {
		return models.UploadPresignResponse{}, err
	}

	return uploadPresignResponse, nil
}

func UploadFile(filePath string, url string, progressCh chan float64) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Получаем размер файла
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	totalSize := fileInfo.Size()

	// Оборачиваем файл в ProgressReader
	progressReader := &models.ProgressReader{
		Reader:     file,
		Total:      totalSize,
		ProgressCh: progressCh,
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest(http.MethodPut, url, progressReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("X-Amz-ACL", "public-read")
	req.ContentLength = totalSize

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file, status: %s", resp.Status)
	}

	return nil
}
