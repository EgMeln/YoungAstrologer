// Package repository provides interfaces and implementations for interacting with data storage.
package repository

import (
	"database/sql"

	"github.com/EgMeln/YoungAstrologer/internal/model"
)

// ImageManager defines the interface for managing images.
type ImageManager interface {
	Create(image *model.Image) error
	GetByDate(date string) (*model.Image, error)
	GetAll() ([]*model.Image, error)
}

// NewImageManager returns a new instance of ImageManager.
func NewImageManager(db *sql.DB) ImageManager {
	return &imageManager{
		db: db,
	}
}

type imageManager struct {
	db *sql.DB
}

// Create inserts a new image into the images table.
func (im *imageManager) Create(image *model.Image) error {
	query := `INSERT INTO images (id, date, explanation, media_type, title, data) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (date) DO NOTHING`

	tx, err := im.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, image.ID, image.Date, image.Explanation, image.MediaType, image.Title, image.Data)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// GetByDate retrieves an image from the images table by the specified date.
func (im *imageManager) GetByDate(date string) (*model.Image, error) {
	query := `SELECT id, date, explanation, media_type, title, data FROM images WHERE date = $1`

	var image model.Image
	tx, err := im.db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(query, date).Scan(&image.ID, &image.Date, &image.Explanation, &image.MediaType, &image.Title, &image.Data)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetAll retrieves all images from the images table.
func (im *imageManager) GetAll() ([]*model.Image, error) {
	query := `SELECT id, date, explanation, media_type, title, data FROM images`

	var images []*model.Image
	tx, err := im.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(query)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var image model.Image

		err := rows.Scan(&image.ID, &image.Date, &image.Explanation, &image.MediaType, &image.Title, &image.Data)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		images = append(images, &image)
	}
	if err := rows.Err(); err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return images, nil
}
