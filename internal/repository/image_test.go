package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/EgMeln/YoungAstrologer/internal/model"
)

func TestImageManager_Create(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE images CASCADE")
		require.NoError(t, err)
	}()

	image := &model.Image{
		ID:          uuid.New(),
		Date:        "2024-05-18",
		Title:       "A Beautiful Nebula",
		Explanation: "This is an explanation of the beautiful nebula.",
		MediaType:   "image",
		Data:        []byte{0x89, 0x50, 0x4E, 0x47},
	}

	err := imageRep.Create(image)
	require.NoError(t, err)

}
func TestImageManager_GetByDate(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE images CASCADE")
		require.NoError(t, err)
	}()

	image := &model.Image{
		ID:          uuid.New(),
		Date:        "2024-05-18",
		Title:       "A Beautiful Nebula",
		Explanation: "This is an explanation of the beautiful nebula.",
		MediaType:   "image",
		Data:        []byte{0x89, 0x50, 0x4E, 0x47},
	}

	err := imageRep.Create(image)
	require.NoError(t, err)

	retrievedImage, err := imageRep.GetByDate(image.Date)
	require.NoError(t, err)

	require.Equal(t, image, retrievedImage)
}

func TestImageManager_GetAll(t *testing.T) {
	defer func() {
		_, err := db.Exec("TRUNCATE TABLE images CASCADE")
		require.NoError(t, err)
	}()

	image1 := &model.Image{
		ID:          uuid.New(),
		Date:        "2024-05-18",
		Explanation: "This is an explanation of the beautiful nebula.",
		MediaType:   "image",
		Title:       "A Beautiful Nebula",
		Data:        []byte{0x89, 0x50, 0x4E, 0x47},
	}

	image2 := &model.Image{
		ID:          uuid.New(),
		Date:        "2024-05-19",
		Explanation: "This is another explanation of the beautiful nebula.",
		MediaType:   "image",
		Title:       "Another Beautiful Nebula",
		Data:        []byte{0x89, 0x50, 0x4E, 0x47},
	}

	err := imageRep.Create(image1)
	require.NoError(t, err)

	err = imageRep.Create(image2)
	require.NoError(t, err)

	images, err := imageRep.GetAll()
	require.NoError(t, err)

	require.Len(t, images, 2)

	require.Equal(t, []*model.Image{image1, image2}, images)
}
