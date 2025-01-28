// Package ui +build !windows
package ui

import (
	"fmt"
	"log"
	"path/filepath"
	"ss/internal/network"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Start() {
	// Create Fyne application
	myApp := app.New()
	myWindow := myApp.NewWindow("File Upload GUI")

	label := widget.NewLabel("Select a file to upload:")
	fileLabel := widget.NewLabel("No file selected")
	responseLabel := widget.NewMultiLineEntry()
	responseLabel.SetPlaceHolder("Response will appear here...")

	// Button to select file
	selectButton := widget.NewButton("Select File", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				log.Println("Error selecting file:", err)
				return
			}
			if reader == nil {
				// User canceled file selection
				return
			}
			fileLabel.SetText(reader.URI().Path())
			defer reader.Close()
		}, myWindow)
		fileDialog.Show()
	})

	// Button to upload file
	uploadButton := widget.NewButton("Upload File", func() {
		filePath := fileLabel.Text
		if filePath == "No file selected" {
			responseLabel.SetText("Please select a file first.")
			return
		}

		// Perform upload logic
		apiKey, err := network.SignInImmediate()
		if err != nil {
			responseLabel.SetText(fmt.Sprintf("Error signing in: %v", err))
			return
		}

		fileName := filepath.Base(filePath)
		presignResponse, err := network.UploadPresign(fileName, apiKey)
		if err != nil {
			responseLabel.SetText(fmt.Sprintf("Error getting presigned URL: %v", err))
			return
		}

		// Создаем канал для получения прогресса
		progressCh := make(chan float64)

		// Запускаем горутину для отслеживания прогресса
		go func() {
			for progress := range progressCh {
				responseLabel.SetText(fmt.Sprintf("Uploading... %.2f%%", progress))
			}
		}()

		// Выполняем загрузку файла
		err = network.UploadFile(filePath, presignResponse.URLToUpload, progressCh)
		if err != nil {
			responseLabel.SetText(fmt.Sprintf("Error uploading file: %v", err))
			return
		}

		// Закрываем канал прогресса после завершения
		close(progressCh)

		responseLabel.SetText(fmt.Sprintf("File uploaded successfully:\n%v", presignResponse.URLToLoad))
	})

	myWindow.SetContent(container.NewVBox(
		label,
		fileLabel,
		selectButton,
		uploadButton,
		responseLabel,
	))

	myWindow.Resize(fyne.NewSize(800, 500))
	myWindow.ShowAndRun()
}
