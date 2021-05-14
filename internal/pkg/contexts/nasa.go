package contexts

import (
	"context"
	"net/http"
	"time"

	"github.com/serbi/nasa_apod_collector/internal/pkg/settings"
)

type NasaContext struct {
	context.Context
	context.CancelFunc
	Client *http.Client
}

func NewNasaContext() *NasaContext {
	ctx, cancel := context.WithCancel(context.Background())
	client := &http.Client{Timeout: time.Second * settings.ApodTimeoutInSeconds}

	return &NasaContext{
		Context:    ctx,
		CancelFunc: cancel,
		Client:     client,
	}
}
