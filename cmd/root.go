package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	userLicense string
	rootCmd     = &cobra.Command{}

	// flags
	esURL        string
	enableColour bool
)

// Execute is respomnsible for executing the viper command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&esURL, esURLFlag, esURLDefault, "url for elasticsearch")
	rootCmd.PersistentFlags().BoolVar(&enableColour, colourFlag, enableColourDefault, "Enable/Disable Colour.")
}
