FROM golang:1.17 as build-env

WORKDIR /app

# Caches the dependencies so future builds are faster
COPY ../../go.mod ./
COPY ../../go.sum ./
RUN go mod download

COPY ../.. ./
RUN CGO_ENABLED=0 go build ./cmd/urlshortener_event_sender/

FROM gcr.io/distroless/base
COPY --from=build-env /app/urlshortener_event_sender /bin/urlshortener_event_sender
EXPOSE 8080
CMD ["/bin/urlshortener_event_sender"]
