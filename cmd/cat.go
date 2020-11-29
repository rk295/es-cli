package cmd

import (
	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "cat api commands",
	Long:  "cat api commands",
}

func init() {
	rootCmd.AddCommand(catCmd)
}
