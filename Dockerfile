# syntax=docker/dockerfile:1
FROM golang:1.17-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build ./cmd/urlshortener/
EXPOSE 8080
CMD [ "./urlshortener" ]

FROM golang:1.17-alpine
WORKDIR /app
COPY --from=0 /app ./
EXPOSE 8080
CMD [ "./urlshortener "]