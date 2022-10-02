FROM golang:1.18

ENV SNACK=BATTLE
ENV PORT=8000

WORKDIR /usr/src/app

# Dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download && \
    go mod verify

# Build
COPY cmd ./cmd
COPY internal/ ./internal
RUN go test ./... && \
    go build -v -o /usr/local/bin/snake ./cmd/snake

CMD ["/usr/local/bin/snake"]
