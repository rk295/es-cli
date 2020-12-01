package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

var (
	shardsCmd = &cobra.Command{
		Use:           "shards",
		Short:         "Lists the shards in the cluster.",
		Long:          "Lists the shards in the cluster. Supports sorting and changing the byte unit to use.",
		RunE:          esShards,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// flags
	sortField  string
	byteFormat string
)

const (
	defaultByteFormat = "mb"
	defaultSortField  = "index"
)

func init() {
	shardsCmd.Flags().StringVarP(&sortField, sortFlag, "s", defaultSortField, "Field to sort by, possible to list multiple comma separated See https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-shards.html for full list of fields")
	shardsCmd.Flags().StringVarP(&byteFormat, byteFlag, "b", defaultByteFormat, `Byte unit to use. Valid values are: "b", "k", "kb", "m", "mb", "g", "gb", "t", "tb", "p" or "pb"`)
	catCmd.AddCommand(shardsCmd)
}

func esShards(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		return err
	}
	defer client.Stop()

	shardsSVC := elastic.NewCatShardsService(client)

	shardList, err := shardsSVC.Bytes(byteFormat).Sort(sortField).Do(ctx)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.Render()
	t.SetStyle(tableStyle)

	t.AppendHeader(table.Row{
		"Index",
		"Shard",
		"State",
		"Docs",
		"Store",
		"IP",
		"Node",
	})

	t.SetCaption("%s_cat/shards?format=json&pretty&bytes=%s&s=%s", esURL, byteFormat, sortField)

	for _, s := range shardList {

		state := s.State
		if enableColour {
			switch state = s.State; state {
			case "UNASSIGNED":
				state = color.RedString(s.State)
			case "STARTED":
				state = color.GreenString(s.State)
			}
		}

		t.AppendRow([]interface{}{
			s.Index,
			s.Shard,
			state,
			s.Docs,
			s.Store,
			s.IP,
			s.Node,
		})
	}
	fmt.Println(t.Render())
	return nil
}
