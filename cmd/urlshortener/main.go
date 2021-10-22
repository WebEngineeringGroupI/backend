package main

import (
	"log"
	"net/http"
)

func main() {
	factory := newFactory()

	log.Fatal(http.ListenAndServe(":8080", factory.NewHTTPRouter()))
}
