package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/EgMeln/YoungAstrologer/internal/model"
	"github.com/EgMeln/YoungAstrologer/internal/service"
)

// APODHandler handles HTTP requests related to Astronomy Picture of the Day (APOD).
type APODHandler struct {
	imageService service.ImageService
	client       *http.Client
}

// NewAPODHandler creates a new instance of APODHandler with the provided imageService.
func NewAPODHandler(imageService service.ImageService, client *http.Client) *APODHandler {
	return &APODHandler{
		imageService: imageService,
		client:       client,
	}
}

// FetchAPOD fetches Astronomy Picture of the Day (APOD) data from NASA API using the provided apiKey.
func (ah *APODHandler) FetchAPOD(apiKey string) (*model.APOD, error) {
	url := fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", apiKey)
	resp, err := ah.client.Get(url)
	if err != nil {
		log.Errorf("Error fetching APOD from NASA API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Unexpected status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
		return nil, err
	}

	var apodResponse model.APOD
	if err := json.Unmarshal(body, &apodResponse); err != nil {
		log.Errorf("Error unmarshalling response: %v", err)
		return nil, err
	}

	return &apodResponse, nil
}

// SaveImage fetches an image from the provided URL and saves it with the given date using the image service.
func (ah *APODHandler) SaveImage(apod *model.APOD) error {
	resp, err := ah.client.Get(apod.URL)
	if err != nil {
		log.Errorf("Error fetching image: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Unexpected status code: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading image data: %v", err)
		return err
	}
	return ah.imageService.Save(&model.Image{
		Date:        apod.Date,
		Explanation: apod.Explanation,
		MediaType:   apod.MediaType,
		Title:       apod.Title,
		Data:        imgData,
	})
}
