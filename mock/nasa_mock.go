package mock

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/serbi/nasa_apod_collector/internal/pkg/nasa"
)

type APODProviderMock struct{}

func (apm *APODProviderMock) GetPicture(ctx context.Context, client *http.Client, picTime time.Time) (image *nasa.Image, err error) {
	if picTime.IsZero() {
		return nil, errors.New("something went terribly wrong")
	}
	currentTime := time.Now()
	image = &nasa.Image{
		Date:        currentTime.Format("2006-01-02"),
		Title:       "Dummy title",
		URL:         "https://apod.nasa.gov/apod/image/1607/dummy_image.jpg",
		HDURL:       "https://apod.nasa.gov/apod/image/1607/dummy_image_HD.jpg",
		Explanation: "This is a dummy description",
		ApodDate:    currentTime,
	}
	return
}
