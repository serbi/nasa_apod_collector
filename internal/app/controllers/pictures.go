package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/serbi/nasa_apod_collector/internal/pkg/contexts"

	"github.com/serbi/nasa_apod_collector/internal/pkg/logger"
	"github.com/serbi/nasa_apod_collector/internal/pkg/models"
	"github.com/serbi/nasa_apod_collector/internal/pkg/nasa"
	"github.com/serbi/nasa_apod_collector/internal/pkg/settings"
)

const ComponentPicturesController = "Component Pictures Controller"

func GetPictures(params *models.PicturesGetParams) (res *models.ApiResponse) {
	if ok := validateParams(params); !ok {
		res = &models.ApiResponse{
			Status: http.StatusUnprocessableEntity,
			Body: &models.ApiResponseBody{
				ApiError: models.NewApiError(models.ErrorMsgPublic422),
			},
		}
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentPicturesController,
			Message:   "validateParams returned false",
		})
		return
	}

	timeStart, timeEnd, err := formatParams(params)
	if err != nil {
		res = &models.ApiResponse{
			Status: http.StatusUnprocessableEntity,
			Body: &models.ApiResponseBody{
				ApiError: models.NewApiError(models.ErrorMsgPublic422),
			},
		}
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentPicturesController,
			Message:   "formatParams: " + err.Error(),
		})
		return
	}

	if ok := timeStart.Before(timeEnd); !ok && timeStart != timeEnd {
		res = &models.ApiResponse{
			Status: http.StatusUnprocessableEntity,
			Body: &models.ApiResponseBody{
				ApiError: models.NewApiError(models.ErrorMsgPublic422),
			},
		}
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentPicturesController,
			Message:   "start_date is after to end_date",
		})
		return
	}

	datesToFetch := getDatesToFetch(timeStart, timeEnd)

	logger.PrintLog(&models.LogMessage{
		Reference: ComponentPicturesController,
		Message:   fmt.Sprintf("fetching images for %d date(s)...", len(datesToFetch)),
	})
	startTimer := time.Now()
	APOD := nasa.NewAPODProvider()
	images, err := fetchAPODImages(APOD, datesToFetch)
	if err != nil {
		res = &models.ApiResponse{
			Status: http.StatusInternalServerError,
			Body: &models.ApiResponseBody{
				ApiError: models.NewApiError(models.ErrorMsgPublic500),
			},
		}
		logger.PrintLogError(&models.LogMessage{
			Reference: ComponentPicturesController,
			Message:   "nasa.ApodImage: " + err.Error(),
		})
		return
	}
	logger.PrintLog(&models.LogMessage{
		Reference: ComponentPicturesController,
		Message:   fmt.Sprintf("fetching successful, execution time: %+v", time.Since(startTimer)),
	})

	imagesUrls := collectURLFields(images)
	res = &models.ApiResponse{
		Status: http.StatusOK,
		Body: &models.ApiResponseBody{
			ApiUrls: &models.ApiUrls{
				Urls: imagesUrls,
			},
		},
	}
	return
}

func validateParams(params *models.PicturesGetParams) (ok bool) {
	dateRe := regexp.MustCompile(`^\d{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`)
	ok = dateRe.MatchString(params.StartDate)
	if !ok {
		return
	}
	ok = dateRe.MatchString(params.EndDate)
	if !ok {
		return
	}
	return true
}

func formatParams(params *models.PicturesGetParams) (timeStart time.Time, timeEnd time.Time, err error) {
	timeStart, err = time.Parse("2006-01-02", params.StartDate)
	if err != nil {
		return
	}
	timeEnd, err = time.Parse("2006-01-02", params.EndDate)
	if err != nil {
		return
	}
	return
}

func getDatesToFetch(timeStart time.Time, timeEnd time.Time) (dates []time.Time) {
	daysBetween := int(timeEnd.Sub(timeStart).Hours() / 24)
	singleDate := timeStart
	for i := 0; i <= daysBetween; i++ {
		dates = append(dates, singleDate)
		singleDate = singleDate.AddDate(0, 0, 1)
	}
	return
}

func fetchAPODImages(APOD nasa.APODPictureGetter, datesToFetch []time.Time) (images []*nasa.Image, err error) {
	nc := contexts.NewNasaContext()
	defer nc.CancelFunc()
	singleImage := &nasa.Image{}

	concurrentGoroutines := make(chan struct{}, settings.ConcurrentRequests)
	for i := 0; i < settings.ConcurrentRequests; i++ {
		concurrentGoroutines <- struct{}{}
	}
	done := make(chan bool)
	errs := make(chan error, 1)

	waitForAllJobs := make(chan bool)
	go func() {
		for i := 0; i < len(datesToFetch); i++ {
			<-done
			concurrentGoroutines <- struct{}{}
		}
		waitForAllJobs <- true
	}()
	for _, d := range datesToFetch {
		<-concurrentGoroutines
		go func(date time.Time) {
			singleImage, err = APOD.GetPicture(nc.Context, nc.Client, date)
			if err != nil {
				sendAsyncError(errs, err)
				nc.CancelFunc()
				waitForAllJobs <- true
			}
			images = append(images, singleImage)
			done <- true
		}(d)
	}
	<-waitForAllJobs
	if nc.Context.Err() != nil {
		return nil, <-errs
	}
	return
}

func collectURLFields(images []*nasa.Image) (urls []string) {
	for _, image := range images {
		urls = append(urls, image.URL)
	}
	return
}

func sendAsyncError(ch chan error, err error) {
	select {
	case ch <- err:
	default:
	}
}
