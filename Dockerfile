# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

# copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY .secrets ./.secrets

RUN CGO_ENABLED=0 GOOS=linux go build -o /oolio ./cmd/server/main.go

# copy static files
COPY static ./static
COPY migrations ./migrations

EXPOSE 8080

CMD ["/oolio"]
