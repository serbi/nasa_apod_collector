package settings

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	ContentTypeHeader = "Content-Type"
	JsonMime          = "application/json"
)

var (
	AppPort               = getDotEnv("PORT", "8080")
	ConcurrentRequests, _ = strconv.Atoi(getDotEnv("CONCURRENT_REQUESTS", "5"))
)

func getDotEnv(key, fallback string) string {
	_ = godotenv.Load()
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
