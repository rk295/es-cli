package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "es-cli",
		Short: "A simple Elasticsearch command line interface",
	}

	// flags
	esURL          string
	enableColour   bool
	markdownOutput bool
)

// Execute is respomnsible for executing the viper command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&esURL, esURLFlag, esURLDefault, "url for elasticsearch")
	rootCmd.PersistentFlags().BoolVar(&enableColour, colourFlag, enableColourDefault, "Enable/Disable Colour.")
	rootCmd.PersistentFlags().BoolVarP(&markdownOutput, markdownFlag, markdownShortFlag, false, "Produce Markdown output")
}
