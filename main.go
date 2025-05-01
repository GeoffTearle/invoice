package main

import (
	_ "embed"
	"log"

	"github.com/maaslalani/invoice/cmd"
)

func main() {
	err := cmd.Root().Execute()
	if err != nil {
		log.Fatal(err)
	}
}
