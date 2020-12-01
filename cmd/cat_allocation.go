package cmd

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

var (
	allocationCmd = &cobra.Command{
		Use:           "allocation",
		Short:         "cat cluster allocations",
		Long:          "cat cluster allocations",
		RunE:          esAllocation,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// flags
	allocSortField string
)

func init() {
	allocationCmd.Flags().StringVarP(&allocSortField, sortFlag, "s", "node", "Field to sort by, possible to list multiple comma separated See https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-allocation.html for full list of fields")
	allocationCmd.Flags().StringVarP(&byteFormat, byteFlag, "b", defaultByteFormat, `Byte unit to use. Valid values are: "b", "k", "kb", "m", "mb", "g", "gb", "t", "tb", "p" or "pb"`)
	catCmd.AddCommand(allocationCmd)
}

func esAllocation(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		return err
	}
	defer client.Stop()

	allocService := elastic.NewCatAllocationService(client)

	// TODO: Add NodeID() to this chain
	allocList, err := allocService.Bytes(byteFormat).Sort(allocSortField).Do(ctx)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.Render()
	t.SetStyle(tableStyle)

	t.AppendHeader(table.Row{
		"Shards",
		"Number of Indices",
		fmt.Sprintf("Disk Used (%s)", byteFormat),
		fmt.Sprintf("Disk Avail (%s)", byteFormat),
		fmt.Sprintf("Disk Total (%s)", byteFormat),
		"Disk Percent",
		"Host",
		"IP",
		"Node",
	})

	t.SetCaption("%s_cat/allocation?format=json&pretty&bytes=%s&s=%s", esURL, byteFormat, allocSortField)

	for _, a := range allocList {

		node := a.Node
		if enableColour {
			switch node = a.Node; node {
			case "UNASSIGNED":
				node = color.RedString(a.Node)
			case "STARTED":
				node = color.GreenString(a.Node)
			}
		}

		t.AppendRow([]interface{}{
			a.Shards,
			a.DiskIndices,
			a.DiskUsed,
			a.DiskAvail,
			a.DiskTotal,
			a.DiskPercent,
			a.Host,
			a.IP,
			node,
		})
	}

	fmt.Println(t.Render())
	return nil
}
