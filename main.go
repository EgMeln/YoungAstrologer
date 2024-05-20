package main

import (
	"crypto/tls"
	"database/sql"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	"github.com/EgMeln/YoungAstrologer/internal/handler"
	"github.com/EgMeln/YoungAstrologer/internal/repository"
	"github.com/EgMeln/YoungAstrologer/internal/service"
)

func main() {
	log.SetLevel(log.InfoLevel)

	logFile, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	nasaAPIKey := os.Getenv("YA_NASA_API_KEY")
	if nasaAPIKey == "" {
		log.Fatal("YA_NASA_API_KEY environment variable is required")
	}

	postgresURL := os.Getenv("YA_POSTGRES_URL")
	if postgresURL == "" {
		log.Fatal("YA_POSTGRES_URL environment variable is required")
	}

	serverPort := os.Getenv("YA_SERVER_PORT")
	if serverPort == "" {
		log.Fatal("YA_SERVER_PORT environment variable is required")
	}

	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v\n", err)
	}
	defer db.Close()

	imageRepo := repository.NewImageManager(db)
	imageSvc := service.NewImageService(imageRepo)
	imageHandler := handler.NewImageHandler(imageSvc)
	apodHandler := handler.NewAPODHandler(imageSvc, client)

	http.HandleFunc("/images", imageHandler.GetAll)
	http.HandleFunc("/images/date", imageHandler.GetByDate)

	done := make(chan bool)
	go startDailyTask(nasaAPIKey, done, apodHandler)

	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}

func startDailyTask(apiKey string, done chan bool, apodHandler *handler.APODHandler) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			ticker = time.NewTicker(24 * time.Hour)

			log.Info("Fetching APOD...")
			apod, err := apodHandler.FetchAPOD(apiKey)
			if err != nil {
				log.Errorf("Error fetching APOD: %v\n", err)
			}

			log.Infof("Title: %s\nDate: %s\nExplanation: %s\nURL: %s\n", apod.Title, apod.Date, apod.Explanation, apod.URL)
			if err := apodHandler.SaveImage(apod); err != nil {
				log.Errorf("Error saving image: %v\n", err)
			} else {
				log.Infof("Image saved for date %s\n", apod.Date)
			}
		}
	}
}
