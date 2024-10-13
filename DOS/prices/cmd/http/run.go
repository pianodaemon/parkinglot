package main

import (
	"log"
	"syscall"

	platform "blaucorp.com/prices/internal/service"
)

const name = "PricesManagerSidecar"

func main() {

	if err := platform.Engage(); err != nil {
		log.Fatalf("%s service struggles with (%v)", name, err)
	}

	syscall.Exit(0)
}
