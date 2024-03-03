FROM golang:1.21.3

ENV GO111MODULE=on

EXPOSE 8080

WORKDIR /app
COPY .. .

RUN go mod download

RUN go build -o imagesetservice cmd/imagesetservice/main.go

CMD ["./imagesetservice"]
