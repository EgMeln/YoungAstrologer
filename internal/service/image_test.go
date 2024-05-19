package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/EgMeln/YoungAstrologer/internal/model"
)

type mockImageManager struct {
	CreateFunc    func(image *model.Image) error
	GetByDateFunc func(date string) (*model.Image, error)
	GetAllFunc    func() ([]*model.Image, error)
}

func (m *mockImageManager) Create(image *model.Image) error {
	return m.CreateFunc(image)
}

func (m *mockImageManager) GetByDate(date string) (*model.Image, error) {
	return m.GetByDateFunc(date)
}

func (m *mockImageManager) GetAll() ([]*model.Image, error) {
	return m.GetAllFunc()
}
func TestImageService_Save(t *testing.T) {
	t.Parallel()

	mockManager := &mockImageManager{
		CreateFunc: func(image *model.Image) error {
			return nil
		},
	}

	imageSvc := NewImageService(mockManager)

	image := &model.Image{
		Date:        "2024-05-18",
		Title:       "A Beautiful Nebula",
		Explanation: "This is an explanation of the beautiful nebula.",
		MediaType:   "image",
		Data:        []byte{0x89, 0x50, 0x4E, 0x47},
	}

	err := imageSvc.Save(image)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, image.ID)
}

func TestImageService_GetByDate(t *testing.T) {
	t.Parallel()

	mockManager := &mockImageManager{
		GetByDateFunc: func(date string) (*model.Image, error) {
			if date == "2024-05-18" {
				return &model.Image{
					ID:          uuid.New(),
					Date:        "2024-05-18",
					Title:       "A Beautiful Nebula",
					Explanation: "This is an explanation of the beautiful nebula.",
					MediaType:   "image",
					Data:        []byte{0x89, 0x50, 0x4E, 0x47},
				}, nil
			}
			return nil, nil
		},
	}

	imageSvc := NewImageService(mockManager)

	image, err := imageSvc.GetByDate("2024-05-18")
	require.NoError(t, err)
	require.NotNil(t, image)
	require.Equal(t, "2024-05-18", image.Date)
}

func TestImageService_GetAll(t *testing.T) {
	t.Parallel()

	mockManager := &mockImageManager{
		GetAllFunc: func() ([]*model.Image, error) {
			return []*model.Image{
				{
					ID:          uuid.New(),
					Date:        "2024-05-18",
					Title:       "A Beautiful Nebula",
					Explanation: "This is an explanation of the beautiful nebula.",
					MediaType:   "image",
					Data:        []byte{0x89, 0x50, 0x4E, 0x47},
				},
				{
					ID:          uuid.New(),
					Date:        "2024-05-19",
					Title:       "Another Beautiful Nebula",
					Explanation: "This is an explanation of another beautiful nebula.",
					MediaType:   "image",
					Data:        []byte{0x89, 0x50, 0x4E, 0x47},
				},
			}, nil
		},
	}

	imageSvc := NewImageService(mockManager)

	images, err := imageSvc.GetAll()
	require.NoError(t, err)
	require.Len(t, images, 2)
}
