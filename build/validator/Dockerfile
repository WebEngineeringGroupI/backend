FROM golang:1.17 as build-env

WORKDIR /app

# Caches the dependencies so future builds are faster
COPY ../../go.mod ./
COPY ../../go.sum ./
RUN go mod download

COPY ../.. ./
RUN CGO_ENABLED=0 go build ./cmd/validator/

FROM gcr.io/distroless/base
COPY --from=build-env /app/validator /bin/validator
EXPOSE 8080
CMD ["/bin/validator"]
