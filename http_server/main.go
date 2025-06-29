package main

import (
	"http_server/api"
	"log"
)

func main() {
	server := api.NewServer("8081")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
