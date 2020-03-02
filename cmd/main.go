package main

import (
	"log"
	"os"

	"github.com/jenkins-x-labs/gsm-controller/pkg"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("no argument found for GCP project id")
	}

	err := pkg.Foo(os.Args[1])
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}
