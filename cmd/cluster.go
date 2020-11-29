package cmd

import (
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cat api commands",
	Long:  "clustercluster api commands",
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}
