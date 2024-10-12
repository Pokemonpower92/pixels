FROM golang:1.22
WORKDIR /app

EXPOSE 8080

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

CMD ["./bin/imagesetparser"]
