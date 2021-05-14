FROM golang:latest

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o nasa_apod_collector cmd/nasa_apod_collector/main.go

CMD [ "/app/nasa_apod_collector" ]