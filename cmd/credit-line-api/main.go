package main

import (
	"log"

	"credit-line/cmd/credit-line-api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
