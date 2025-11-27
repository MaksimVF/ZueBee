


FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /tail-service ./cmd/tail

EXPOSE 8080

ENTRYPOINT ["/tail-service"]


