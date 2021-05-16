package nasa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/serbi/nasa_apod_collector/internal/pkg/settings"
)

var apodEndpoint = settings.ApodEndpoint

var nasaApiKey = settings.NasaApiKey

type Image struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	HDURL       string `json:"hdurl"`
	Explanation string `json:"explanation"`

	ApodDate time.Time `json:",omitempty"`
}

func (ni Image) String() string {
	return fmt.Sprintf(`Title: %s
Date: %s
Image: %s
HD Image: %s
About:
%s
`, ni.Title, ni.Date, ni.URL, ni.HDURL, ni.Explanation)
}

type APODPictureGetter interface {
	GetPicture(ctx context.Context, client *http.Client, picTime time.Time) (*Image, error)
}

type APODProvider struct {
	baseURL string
	apiKey  string
}

func NewAPODProvider() *APODProvider {
	return &APODProvider{
		baseURL: apodEndpoint,
		apiKey:  nasaApiKey,
	}
}

func (ap *APODProvider) GetPicture(ctx context.Context, client *http.Client, picTime time.Time) (*Image, error) {
	if err := validateDate(picTime); err != nil {
		return nil, err
	}

	req, err := createAPODRequest(ctx, ap, picTime)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to NASA API, %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	var ni Image
	err = json.Unmarshal(dat, &ni)
	if err != nil {
		return nil, err
	}
	if ni.URL == "" && ni.HDURL == "" {
		return nil, errors.New("NASA APOD API is returned an invalid response, may be down temporarily")
	}
	if t, err := time.Parse("2006-01-02", ni.Date); err == nil {
		ni.ApodDate = t
	}
	return &ni, nil
}

func validateDate(picTime time.Time) error {
	if picTime.Before(time.Date(1995, 6, 16, 0, 0, 0, 0, time.UTC)) ||
		picTime.After(time.Now()) {
		return errors.New("request date must be between Jun 16, 1995 and today")
	}
	return nil
}

func createAPODRequest(ctx context.Context, ap *APODProvider, picTime time.Time) (req *http.Request, err error) {
	date := picTime.Format("2006-01-02")
	u, err := url.Parse(ap.baseURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("api_key", ap.apiKey)
	q.Add("date", date)
	u.RawQuery = q.Encode()
	req, err = http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return
}
