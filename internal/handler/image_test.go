package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/EgMeln/YoungAstrologer/internal/model"
)

type mockImageService struct {
	GetByDateFunc func(date string) (*model.Image, error)
	GetAllFunc    func() ([]*model.Image, error)
	SaveFunc      func(image *model.Image) error
}

func (m *mockImageService) GetByDate(date string) (*model.Image, error) {
	return m.GetByDateFunc(date)
}

func (m *mockImageService) GetAll() ([]*model.Image, error) {
	return m.GetAllFunc()
}
func (m *mockImageService) Save(image *model.Image) error {
	return m.SaveFunc(image)
}

func TestImageHandler_GetByDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		date               string
		getByDateFunc      func(date string) (*model.Image, error)
		expectedStatusCode int
		expectedResponse   *model.Image
	}{
		{
			name: "Success",
			date: "2024-05-18",
			getByDateFunc: func(date string) (*model.Image, error) {
				return &model.Image{
					ID:          uuid.New(),
					Date:        "2024-05-18",
					Title:       "A Beautiful Nebula",
					Explanation: "This is an explanation of the beautiful nebula.",
					MediaType:   "image",
					Data:        []byte{0x89, 0x50, 0x4E, 0x47},
				}, nil
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &model.Image{
				Date:        "2024-05-18",
				Title:       "A Beautiful Nebula",
				Explanation: "This is an explanation of the beautiful nebula.",
				MediaType:   "image",
			},
		},
		{
			name: "ImageNotFound",
			date: "2024-05-19",
			getByDateFunc: func(date string) (*model.Image, error) {
				return nil, nil
			},
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name: "ServiceError",
			date: "2024-05-20",
			getByDateFunc: func(date string) (*model.Image, error) {
				return nil, errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "BadRequest",
			date:               "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageService := &mockImageService{
				GetByDateFunc: tt.getByDateFunc,
			}
			imageHandler := NewImageHandler(imageService)

			req, err := http.NewRequest(http.MethodGet, "/images/date?date="+tt.date, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			imageHandler.GetByDate(recorder, req)

			require.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedResponse != nil {
				var image model.Image
				err = json.NewDecoder(recorder.Body).Decode(&image)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResponse.Date, image.Date)
				require.Equal(t, tt.expectedResponse.Title, image.Title)
				require.Equal(t, tt.expectedResponse.Explanation, image.Explanation)
				require.Equal(t, tt.expectedResponse.MediaType, image.MediaType)
			}
		})
	}
}

func TestImageHandler_GetAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		getAllFunc         func() ([]*model.Image, error)
		expectedStatusCode int
		expectedResponse   []*model.Image
	}{
		{
			name: "Success",
			getAllFunc: func() ([]*model.Image, error) {
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
			expectedStatusCode: http.StatusOK,
			expectedResponse: []*model.Image{
				{
					Date:        "2024-05-18",
					Title:       "A Beautiful Nebula",
					Explanation: "This is an explanation of the beautiful nebula.",
					MediaType:   "image",
				},
				{
					Date:        "2024-05-19",
					Title:       "Another Beautiful Nebula",
					Explanation: "This is an explanation of another beautiful nebula.",
					MediaType:   "image",
				},
			},
		},
		{
			name: "ServiceError",
			getAllFunc: func() ([]*model.Image, error) {
				return nil, errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageService := &mockImageService{
				GetAllFunc: tt.getAllFunc,
			}
			imageHandler := NewImageHandler(imageService)

			req, err := http.NewRequest(http.MethodGet, "/images/all", nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			imageHandler.GetAll(recorder, req)

			require.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedResponse != nil {
				var images []*model.Image
				err = json.NewDecoder(recorder.Body).Decode(&images)
				require.NoError(t, err)
				require.Equal(t, len(tt.expectedResponse), len(images))
				for i, image := range images {
					require.Equal(t, tt.expectedResponse[i].Date, image.Date)
					require.Equal(t, tt.expectedResponse[i].Title, image.Title)
					require.Equal(t, tt.expectedResponse[i].Explanation, image.Explanation)
					require.Equal(t, tt.expectedResponse[i].MediaType, image.MediaType)
				}
			}
		})
	}
}
