package cmd

import (
	"context"
	"fmt"

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

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Configuration", Transformer: prettyJSONTransformer()},
	})

	t.SetCaption("%s_snapshot", esURL)

	for n, v := range status {
		t.AppendRow([]interface{}{
			n,
			v,
		})
	}
	fmt.Println(render(t))
	return nil

}
