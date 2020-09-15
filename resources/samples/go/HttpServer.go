package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const port = 3090

func hello(writer http.ResponseWriter, req *http.Request) {
	if fileInfo, err := os.Stat("health.check"); err != nil {
		log.Print("file health.check not found, returning 404")
		writer.WriteHeader(404)
	} else {
		log.Printf("file %s found, returning hello", fileInfo.Name())
		fmt.Fprintf(writer, "hello\n")
	}
}

func main() {
	http.HandleFunc("/", hello)
	log.Printf("starting up http server on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
