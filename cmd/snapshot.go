package cmd

import (
	"context"
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	"github.com/jedib0t/go-pretty/v6/table"
	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

var (
	snapshotCmd = &cobra.Command{
		Use:           "snapshot",
		Short:         "snapshot repositories",
		Long:          "snapshot repositories",
		RunE:          esSnapshot,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func init() {
	rootCmd.AddCommand(snapshotCmd)
}

func esSnapshot(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		return err
	}
	defer client.Stop()

	repositorySVC := elastic.NewSnapshotGetRepositoryService(client)

	status, err := repositorySVC.Do(ctx)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.Render()
	t.SetStyle(tableStyle)

	t.AppendHeader(table.Row{
		"Name",
		"Configuration",
	})

	t.SetCaption("%s_snapshot", esURL)

	for n, v := range status {
		s, err := prettyjson.Marshal(v)
		if err != nil {
			return err
		}
		settings := string(s)

		t.AppendRow([]interface{}{
			n,
			settings,
		})
	}
	fmt.Println(render(t))
	return nil

}
