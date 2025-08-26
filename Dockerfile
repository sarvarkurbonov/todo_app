FROM  golang:alpine AS builder
WORKDIR /application

COPY go.mod go.sum ./

RUN apk add --no-cache git ca-certificates

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o /app/server ./cmd/main.go

CMD ["/app/server"]