FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /webdav


EXPOSE 8080
EXPOSE 8443

ENTRYPOINT ["/webdav"]