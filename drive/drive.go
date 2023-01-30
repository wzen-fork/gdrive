package drive

import (
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"net/http"
)

type Drive struct {
	service *drive.Service
}

func New(client *http.Client) (*Drive, error) {
	service, err := drive.New(client)
	if err != nil {
		return nil, err
	}

	return &Drive{service}, nil
}

func NewWithAPIKey(apiKey string) (*Drive, error) {
	service, err := drive.NewService(context.Background(),option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &Drive{service}, nil
}

