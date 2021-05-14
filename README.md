# nasa_apod_collector
This repository was a part of the recruitment task to show off my Golang skills.
The task was to create a micro-service to collect URLs for given date range from NASA's Astronomy Picture of the Day.
You can check it out here: https://apod.nasa.gov/apod/astropix.html


## Environment variables configuration

To configure environment variables, please create `.env` file in the directory as shown in the
`.env.example` file.

## Available Scripts
In the project directory, you can run:

| Command | Action |
| ------- | ----------- |
| `make build` | Builds the app locally |
| `make run` | Runs the app locally |
| `make test` | Runs golang tests locally |
|  |  |
| `make docker-build` | Runs the docker build command |
| `make docker-run` | Runs the app inside docker container |
| `make docker-test` | Runs go test in the docker container |

By default the app can be accessed on [http://localhost:8000](http://localhost:8000).

## How to
A single endpoint is exposed under the following path:  

`GET` `/pictures?start_date=2020-01-04&end_date=2020-02-05`  

Since the NASA API publishes one image per day the `start_date` and `end_date` parameters define
range of pictures to be processed. The response from the endpoint is a JSON message
containing all the urls in the following format:  

`{“urls”: ["https://apod.nasa.gov/apod/image/2008/AlienThrone_Zajac_3807.jpg", ...]}`  

In case of an error the server is asked to return an appropriate HTTP status code and a descriptive JSON
message in the following format:  

`{“error”: “error message”}`  

- As the provided date range in a single request might be broad, the NASA API is queried
concurrently. However, in order not to be recognized as a malicious user, a limit of concurrent
requests to this external API (configurable by `CONCURRENT_REQUESTS` var) was added. This limit should never be exceeded
regardless of how many concurrent requests is the url-collector receiving.  

---

### My requirements:
- There should be some unit tests - no need to test everything, pick just one component and test it
thoroughly.
- `start_date` and `end_date` parameters should have some kind of validation - think of possible
corner cases (for example: start_date should be earlier than end_date )
- The application should be configurable via environment variables. I should be able to provide
following variables on application startup:
    1. api_key used for NASA API requests (variable name `API_KEY`, default: `DEMO_KEY`)
    2. Limit of concurrent requests to NASA API (variable name `CONCURRENT_REQUESTS`, default: 5)
    3. A port the server is running on (variable name `PORT`, default: `8080`)
- The project should be properly structured to allow the use of some other images providers.
- I should provide a Dockerfile that can be used to build and run the application without having the Go
toolchain installed locally.
