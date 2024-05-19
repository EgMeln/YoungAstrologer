package main

import (
	"crypto/tls"
	"database/sql"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/EgMeln/YoungAstrologer/internal/config"
	"github.com/EgMeln/YoungAstrologer/internal/handler"
	"github.com/EgMeln/YoungAstrologer/internal/repository"
	"github.com/EgMeln/YoungAstrologer/internal/service"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	logrus.SetOutput(logFile)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("Error loading config: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.PostgresURL)
	if err != nil {
		logrus.Fatalf("Error connecting to the database: %v\n", err)
	}
	defer db.Close()

	imageRepo := repository.NewImageManager(db)
	imageSvc := service.NewImageService(imageRepo)
	imageHandler := handler.NewImageHandler(imageSvc)
	apodHandler := handler.NewAPODHandler(imageSvc, client)

	http.HandleFunc("/images", imageHandler.GetAll)
	http.HandleFunc("/images/date", imageHandler.GetByDate)

	done := make(chan bool)
	go startDailyTask(cfg.NASAAPIKey, done, apodHandler)

	if err := http.ListenAndServe(cfg.ServerPort, nil); err != nil {
		logrus.Fatalf("Error starting server: %v\n", err)
	}
}

func startDailyTask(apiKey string, done chan bool, apodHandler *handler.APODHandler) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			ticker = time.NewTicker(24 * time.Hour)

			logrus.Info("Fetching APOD...")

			apod, err := apodHandler.FetchAPOD(apiKey)
			if err != nil {
				logrus.Errorf("Error fetching APOD: %v\n", err)
			}

			logrus.Infof("Title: %s\nDate: %s\nExplanation: %s\nURL: %s\n", apod.Title, apod.Date, apod.Explanation, apod.URL)
			if err := apodHandler.SaveImage(apod); err != nil {
				logrus.Errorf("Error saving image: %v\n", err)
			} else {
				logrus.Infof("Image saved for date %s\n", apod.Date)
			}
		}
	}
}
