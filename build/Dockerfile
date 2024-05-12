FROM golang:1.21.3
WORKDIR /app

EXPOSE 8080

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o imagesetservice cmd/imagesetservice/main.go

CMD ["./imagesetservice"]
