FROM golang:1.18

ENV PORT=8000
ENV SNAKE=LATEST

WORKDIR /usr/src/app

# Dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download && \
    go mod verify

# Build
COPY cmd ./cmd
COPY snacks ./snacks
RUN go test ./... && \
    go build -v -o /usr/local/bin/app ./cmd

CMD ["app"]
