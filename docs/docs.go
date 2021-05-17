package main

import (
	"log"

	"github.com/rk295/es-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	rootCMD := cmd.NewESCliCommand()
	err := doc.GenMarkdownTree(rootCMD, "./")
	if err != nil {
		log.Fatal(err)
	}
}
