package main

import (
	"os"

	"github.com/REPLACE_ME_ORG/REPLACE_ME_APP_NAME/cmd/root"
)

// Entrypoint for the command
func main() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
