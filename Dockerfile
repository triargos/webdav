FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o /webdav
ENV DOCKER_ENABLED="1"
ENTRYPOINT ["/webdav"]