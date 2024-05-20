// Package service provides business logic services for the application.
package service

import (
	"github.com/google/uuid"

	"github.com/EgMeln/YoungAstrologer/internal/model"
	"github.com/EgMeln/YoungAstrologer/internal/repository"
)

// ImageService defines the interface for the image service.
type ImageService interface {
	Save(image *model.Image) error
	GetByDate(date string) (*model.Image, error)
	GetAll() ([]*model.Image, error)
}

// NewImageService returns a new instance of ImageService.
func NewImageService(imageManager repository.ImageManager) ImageService {
	return &imageService{
		imageManager: imageManager,
	}
}

type imageService struct {
	imageManager repository.ImageManager
}

// Save generates a new UUID for the image and stores it in the database.
func (is *imageService) Save(image *model.Image) error {
	image.ID = uuid.New()
	return is.imageManager.Create(image)
}

// GetByDate retrieves an image from the database by the specified date.
func (is *imageService) GetByDate(date string) (*model.Image, error) {
	return is.imageManager.GetByDate(date)
}

// GetAll retrieves all images from the database.
func (is *imageService) GetAll() ([]*model.Image, error) {
	return is.imageManager.GetAll()
}
