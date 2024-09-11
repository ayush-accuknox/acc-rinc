package main

import (
	"log"
	"os"

	"github.com/murtaza-u/rinc/internal/conf"
)

func main() {
	conf, err := conf.New(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
	err = conf.Validate()
	if err != nil {
		log.Fatalf("validating provided config: %s", err.Error())
	}
	log.Println("Reporter IN Cluster")
}
