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
		Use:               "nodes [node-name]",
		Short:             "es nodes",
		Long:              "es cluster nodes",
		RunE:              esNodes,
		SilenceErrors:     true,
		SilenceUsage:      true,
		Args:              cobra.MinimumNArgs(0),
		ValidArgsFunction: esNodeNounCompletion(),
	}
)

func init() {
	rootCmd.AddCommand(nodesCmd)
}

func esNodes(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	node := ""
	if len(args) >= 1 {
		node = args[0]
	}
	status, err := getESNodes(ctx, node)
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

	u := fmt.Sprintf("%s_nodes", esURL)
	if node != "" {
		u = fmt.Sprintf("%s/%s", u, node)
	}
	u = fmt.Sprintf("%s?format=json&pretty", u)
	t.SetCaption(u)

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

func getESNodes(ctx context.Context, node string) (*elastic.NodesInfoResponse, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		return &elastic.NodesInfoResponse{}, err
	}
	defer client.Stop()

	nodes := client.NodesInfo()

	return nodes.NodeId(node).Do(ctx)
}

func getESNodeNames(ctx context.Context) ([]string, error) {

	var nodeNames []string

	nodeInfo, err := getESNodes(ctx, "")
	if err != nil {
		return nodeNames, err
	}

	for _, n := range nodeInfo.Nodes {
		nodeNames = append(nodeNames, n.Name)
	}
	return nodeNames, err
}
