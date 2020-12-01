package cmd

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/table"
	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

var (
	tasksCmd = &cobra.Command{
		Use:           "tasks",
		Short:         "cat tasks",
		Long:          "cat tasks",
		RunE:          esCatTasks,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
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
	})

	t.SetCaption("%s_cat/tasks?format=json&pretty", esURL)

	var tableRows rows

	// TODO: Explore printing the other entries, not just Nodes
	// fmt.Println("Node Failures:", taskList.NodeFailures)
	// fmt.Println("Task Failures:", taskList.TaskFailures)

	for _, node := range taskList.Nodes {
		for _, task := range node.Tasks {
			tableRows = append(tableRows,
				table.Row{
					node.Name,
					fmt.Sprintf("%v:%v", task.Node, task.Id),
					task.Action,
					task.Type,
					time.Unix(0, task.StartTimeInMillis*int64(time.Millisecond)),
					time.Duration(task.RunningTimeInNanos) * time.Nanosecond,
				},
			)
		}
	}

	sort.Slice(tableRows, func(i, j int) bool {
		return tableRows[i][5].(time.Duration) < tableRows[j][5].(time.Duration)
	})

	t.AppendRows(tableRows)
	fmt.Println(t.Render())
	return nil

}
