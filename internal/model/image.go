// Package model contains the data models used in the application.
package model

import (
	"github.com/google/uuid"
)

// Image represents an image entity
type Image struct {
	ID          uuid.UUID
	Date        string
	Explanation string
	MediaType   string
	Title       string
	Data        []byte
}
