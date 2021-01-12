package cmd

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

var (
	tasksCmd = &cobra.Command{
		Use:   "tasks [node-name]",
		Short: "Lists the running tasks within the cluster",
		Long: `Lists the running tasks within the cluster. Sorted by running time.
`,
		RunE:          esCatTasks,
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.MinimumNArgs(0),
	}
)

const (
	redDuration    time.Duration = 30 * time.Second
	yellowDuration time.Duration = 10 * time.Second
)

func init() {
	catCmd.AddCommand(tasksCmd)
}

type rows []table.Row

func esCatTasks(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient(elastic.SetURL(esURL))
	if err != nil {
		return err
	}
	defer client.Stop()

	tasksSVC := elastic.NewTasksListService(client)

	taskList, err := tasksSVC.Do(ctx)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.Render()
	t.SetStyle(tableStyle)

	t.AppendHeader(table.Row{
		"Name",
		"Task ID",
		"Action",
		"Type",
		"Start Time",
		"Running Time",
		"Sort - hidden in output",
	})

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 7, Hidden: true},
	})

	t.SetCaption("%s_cat/tasks?format=json&pretty", esURL)

	var tableRows rows

	// TODO: Explore printing the other entries, not just Nodes
	// fmt.Println("Node Failures:", taskList.NodeFailures)
	// fmt.Println("Task Failures:", taskList.TaskFailures)

	for _, node := range taskList.Nodes {

		if len(args) == 1 && node.Name != args[0] {
			continue
		}

		for _, task := range node.Tasks {

			duration := time.Duration(task.RunningTimeInNanos) * time.Nanosecond
			var prettyDuration string

			if duration > redDuration {
				prettyDuration = color.RedString(fmt.Sprintf("%v", duration))
			} else if duration > yellowDuration {
				prettyDuration = color.YellowString(fmt.Sprintf("%v", duration))
			} else {
				prettyDuration = color.GreenString(fmt.Sprintf("%v", duration))
			}

			tableRows = append(tableRows,
				table.Row{
					node.Name,
					fmt.Sprintf("%v:%v", task.Node, task.Id),
					task.Action,
					task.Type,
					time.Unix(0, task.StartTimeInMillis*int64(time.Millisecond)),
					prettyDuration,
					duration,
				},
			)
		}
	}

	sort.Slice(tableRows, func(i, j int) bool {
		return tableRows[i][6].(time.Duration) < tableRows[j][6].(time.Duration)
	})

	t.AppendRows(tableRows)
	fmt.Println(t.Render())
	return nil

}
