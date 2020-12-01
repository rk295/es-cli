package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:           "health",
	Short:         "health",
	Long:          "health",
	RunE:          catHealth,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	catCmd.AddCommand(healthCmd)
}

func catHealth(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		return err
	}
	defer client.Stop()

	healthService := elastic.NewCatHealthService(client)

	health, err := healthService.Do(ctx)
	if err != nil {
		return err
	}

	var s []byte
	if enableColour {
		s, _ = prettyjson.Marshal(health)
	} else {
		s, _ = json.MarshalIndent(health, "", " ")
	}

	fmt.Println(string(s))
	return nil
}
