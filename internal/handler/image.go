// Package handler provides HTTP handlers for managing images in the image library.
package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/EgMeln/YoungAstrologer/internal/service"
)

// ImageHandler handles HTTP requests related to images.
type ImageHandler struct {
	imageService service.ImageService
}

// NewImageHandler creates a new ImageHandler instance.
func NewImageHandler(imageService service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

// GetByDate handles the HTTP request for retrieving an image by date.
func (ih *ImageHandler) GetByDate(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		log.Warn("Date parameter is missing")
		http.Error(w, "Date parameter is required", http.StatusBadRequest)
		return
	}

	image, err := ih.imageService.GetByDate(date)
	if err != nil {
		log.Errorf("Failed to get image by date: %v", err)
		http.Error(w, "Failed to get image by date", http.StatusInternalServerError)
		return
	}
	if image == nil {
		log.Errorf("Image not found for date: %s", date)
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(image); err != nil {
		log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetAll handles the HTTP request for retrieving all images.
func (ih *ImageHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	images, err := ih.imageService.GetAll()
	if err != nil {
		log.Errorf("Failed to get all images: %v", err)
		http.Error(w, "Failed to get all images", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(images); err != nil {
		log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
