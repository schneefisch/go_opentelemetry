package main

import (
	"github.com/schneefisch/go_opentelemetry/app"
	"net/http"
)

func main() {
	router := app.InitApp()
	_ = http.ListenAndServe(":8080", router)
}
