package cmd

// TODO: Allow specifying a single node

import (
	"context"
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	"github.com/jedib0t/go-pretty/v6/table"
	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

var (
	nodesCmd = &cobra.Command{
		Use:           "nodes",
		Short:         "es nodes",
		Long:          "es cluster nodes",
		RunE:          esNodes,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func init() {
	rootCmd.AddCommand(nodesCmd)
}

func esNodes(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		return err
	}
	defer client.Stop()

	nodes := client.NodesInfo()

	status, err := nodes.Do(ctx)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.Render()
	t.SetStyle(tableStyle)

	t.AppendHeader(table.Row{
		"Name",
		"Version",
		"Attributes",
		"Node",
	})

	t.SetCaption("%s_nodes", esURL)

	for _, v := range status.Nodes {

		var attributes string

		if v.Attributes == nil {
			attributes = ""
		} else {
			a, err := prettyjson.Marshal(v.Attributes)
			if err != nil {
				return err
			}
			attributes = string(a)
		}

		var settings string
		if v.Settings == nil {
			settings = ""
		} else {
			s, err := prettyjson.Marshal(v.Settings)
			if err != nil {
				return err
			}
			settings = string(s)
		}

		t.AppendRow([]interface{}{
			v.Name,
			v.Version,
			attributes,
			settings,
		})
	}
	fmt.Println(render(t))
	return nil

}
