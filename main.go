package main

import (
	"net/http"
	"time"

	"fmt"
	"os"

	"github.com/roman-kulish/wfh/generator"
)

const port = 80

var addr string

func main() {
	mux := *http.NewServeMux()
	mux.Handle(generator.Command, generator.NewHandler())

	server := http.Server{
		Addr:         addr,
		Handler:      &mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 3,
		IdleTimeout:  time.Second * 10,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func init() {
	if listen, ok := os.LookupEnv("APP_PORT"); ok {
		addr = fmt.Sprintf(":%s", listen)
	} else {
		addr = fmt.Sprintf(":%d", port)
	}
}
