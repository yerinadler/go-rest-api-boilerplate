package main

import (
	"fmt"
	"log"

	"github.com/example/go-rest-api-revision/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Kafka.Brokers)
}
