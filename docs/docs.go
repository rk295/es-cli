package main

import (
	"log"

	"github.com/rk295/es-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	foo := cmd.NewESCliCommand()
	err := doc.GenMarkdownTree(foo, "./")
	if err != nil {
		log.Fatal(err)
	}
}
