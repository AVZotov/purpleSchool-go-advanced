package main

import (
	"fmt"
	"log"
	"order_api_auth/internal/config"
	"os"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Application panicked: %v", rec)
			os.Exit(1)
		}
	}()
	const DevFile = "configs.yml"

	cfg, err := config.MustLoadConfig("configs.yml")
	if err != nil {
		//log.Fatalf("Error loading config: %v", err)
		panic(err)
	}
	fmt.Printf("%+v\n", cfg)
}
