package settings

const (
	ApodEndpoint         = "https://api.nasa.gov/planetary/apod"
	ApodTimeoutInSeconds = 20
)

var (
	NasaApiKey = getDotEnv("API_KEY", "DEMO_KEY")
)
