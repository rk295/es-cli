package cmd

// TODO: Allow specifying a single node

import (
	"context"
	"fmt"

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

	status, err := getESNodes(ctx)
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

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Attributes", Transformer: prettyJSONTransformer()},
		{Name: "Node", Transformer: prettyJSONTransformer()},
	})

	t.SetCaption("%s_nodes", esURL)

	for _, v := range status.Nodes {

		t.AppendRow([]interface{}{
			v.Name,
			v.Version,
			v.Attributes,
			v.Settings,
		})
	}
	fmt.Println(render(t))
	return nil

}

func getESNodes(ctx context.Context) (*elastic.NodesInfoResponse, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		return &elastic.NodesInfoResponse{}, err
	}
	defer client.Stop()

	nodes := client.NodesInfo()

	return nodes.Do(ctx)
}

func getESNodeNames(ctx context.Context) ([]string, error) {

	var nodeNames []string

	nodeInfo, err := getESNodes(ctx)
	if err != nil {
		return nodeNames, err
	}

	for _, n := range nodeInfo.Nodes {
		nodeNames = append(nodeNames, n.Name)
	}
	return nodeNames, err
}
