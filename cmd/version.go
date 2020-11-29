package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  "Print version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Some version string...")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
