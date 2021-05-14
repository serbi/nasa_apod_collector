docker-build:
	docker build -t nasa_apod_collector .

docker-run:
	docker run --rm -it -p 8080:8080 nasa_apod_collector

docker-test:
	docker run nasa_apod_collector sh -c "go test ./..."