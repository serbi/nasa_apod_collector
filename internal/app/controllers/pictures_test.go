package controllers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/serbi/nasa_apod_collector/internal/pkg/models"
	"github.com/serbi/nasa_apod_collector/internal/pkg/nasa"
	"github.com/serbi/nasa_apod_collector/mock"
)

func TestValidateParams(t *testing.T) {
	t.Run("should return true for valid date params", func(t *testing.T) {
		dummyParams := &models.PicturesGetParams{
			StartDate: "1993-03-14",
			EndDate:   "1997-06-27",
		}
		ok := validateParams(dummyParams)
		assert.True(t, ok)
	})
	t.Run("should return false if date params are in wrong format", func(t *testing.T) {
		dummyParams := &models.PicturesGetParams{
			StartDate: "03-14-1993",
			EndDate:   "06-27-1997",
		}
		ok := validateParams(dummyParams)
		assert.False(t, ok)

		dummyParams = &models.PicturesGetParams{
			StartDate: "03.14.1993",
			EndDate:   "06.27.1997",
		}
		ok = validateParams(dummyParams)
		assert.False(t, ok)
	})
	t.Run("should return false if date is empty", func(t *testing.T) {
		dummyParams := &models.PicturesGetParams{
			StartDate: "",
			EndDate:   "",
		}
		ok := validateParams(dummyParams)
		assert.False(t, ok)
	})
}

func TestFormatParams(t *testing.T) {
	t.Run("should format date params and return error nil for valid date params", func(t *testing.T) {
		dummyParams := &models.PicturesGetParams{
			StartDate: "1993-03-14",
			EndDate:   "1997-06-27",
		}
		timeStart, timeEnd, err := formatParams(dummyParams)
		assert.IsType(t, time.Time{}, timeStart)
		assert.IsType(t, time.Time{}, timeEnd)
		assert.Nil(t, err)
	})
	t.Run("should return error for invalid date params", func(t *testing.T) {
		dummyParams := &models.PicturesGetParams{
			StartDate: "06-27-1997",
			EndDate:   "dummy",
		}
		_, _, err := formatParams(dummyParams)
		assert.NotNil(t, err)
	})
	t.Run("should return error for empty date params", func(t *testing.T) {
		dummyParams := &models.PicturesGetParams{
			StartDate: "",
			EndDate:   "",
		}
		_, _, err := formatParams(dummyParams)
		assert.NotNil(t, err)
	})
}

func TestGetDatesToFetch(t *testing.T) {
	t.Run("should return a list of dates", func(t *testing.T) {
		timeStart := time.Date(2019, 10, 10, 0, 0, 0, 0, time.UTC)
		timeEnd := time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC)
		dates := getDatesToFetch(timeStart, timeEnd)
		assert.IsType(t, time.Time{}, dates[0])
		assert.Equal(t, 367, len(dates))
	})
}

func TestFetchAPODImages(t *testing.T) {
	nasaContextMock := &mock.APODProviderMock{}
	t.Run("should properly fetch APOD image urls", func(t *testing.T) {
		datesToFetch := []time.Time{
			time.Date(2020, 04, 01, 00, 00, 00, 0000, time.UTC),
			time.Date(2020, 04, 02, 00, 00, 00, 0000, time.UTC),
			time.Date(2020, 04, 03, 00, 00, 00, 0000, time.UTC),
		}
		urls, err := fetchAPODImages(nasaContextMock, datesToFetch)
		assert.Nil(t, err)
		assert.Equal(t, len(datesToFetch), len(urls))
		expectedUrl := "https://apod.nasa.gov/apod/image/1607/dummy_image.jpg"
		assert.Equal(t, expectedUrl, urls[0].URL)
		assert.Equal(t, expectedUrl, urls[1].URL)
		assert.Equal(t, expectedUrl, urls[2].URL)
	})
	t.Run("should fail after first invalid goroutine request", func(t *testing.T) {
		datesToFetch := []time.Time{
			time.Date(2020, 04, 01, 00, 00, 00, 0000, time.UTC),
			time.Date(2020, 04, 02, 00, 00, 00, 0000, time.UTC),
			time.Date(2020, 04, 03, 00, 00, 00, 0000, time.UTC),
			{},
		}
		urls, err := fetchAPODImages(nasaContextMock, datesToFetch)
		assert.NotNil(t, err)
		assert.Equal(t, 0, len(urls))
	})
}

func TestCollectURLFields(t *testing.T) {
	t.Run("should return a list of urls given the images list", func(t *testing.T) {
		images := []*nasa.Image{
			{
				"2020-04-05",
				"Dummy title",
				"https://apod.nasa.gov/apod/image/1607/BeyondEarth_Unknown_960.jpg",
				"https://apod.nasa.gov/apod/image/1607/BeyondEarth_Unknown_3000.jpg",
				"This is dummy description",
				time.Now(),
			},
			{
				"1993-04-05",
				"Other dummy title",
				"https://apod.nasa.gov/apod/image/2004/PotatoPod_Sutton_960.jpg",
				"https://apod.nasa.gov/apod/image/2004/PotatoPod_Sutton_5332.jpg",
				"This is another dummy description",
				time.Now(),
			},
		}
		expectedUrls := []string{
			"https://apod.nasa.gov/apod/image/1607/BeyondEarth_Unknown_960.jpg",
			"https://apod.nasa.gov/apod/image/2004/PotatoPod_Sutton_960.jpg",
		}
		urls := collectURLFields(images)
		assert.Equal(t, expectedUrls, urls)
	})
}
