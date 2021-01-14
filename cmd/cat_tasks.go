package cmd

import (
	"context"
	"fmt"
	"sort"
	"time"

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
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			var nodes []string
			nodes, err := getESNodeNames(context.Background())
			if err != nil {
				return nodes, cobra.ShellCompDirectiveNoFileComp
			}
			return nodes, cobra.ShellCompDirectiveNoFileComp
		},
	}

	detailedFlag bool
)

const (
	redDuration    time.Duration = 30 * time.Second
	yellowDuration time.Duration = 10 * time.Second
)

func init() {
	tasksCmd.Flags().BoolVarP(&detailedFlag, "detailed", "d", false, "Include detailed output")
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

	taskList, err := tasksSVC.Detailed(detailedFlag).Do(ctx)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.Render()
	t.SetStyle(tableStyle)
	h := table.Row{
		"Name",
		"Task ID",
		"Action",
		"Type",
		"Start Time",
		"Running Time",
		"Sort - hidden in output",
	}

	if detailedFlag {
		h = append(h, "Description")
	}

	t.AppendHeader(h)

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 5, Transformer: timeInMSTransformer()},
		{Number: 6, Transformer: durationTransformer()},
		{Number: 7, Hidden: true},
	})

	url := fmt.Sprintf("%s_cat/tasks?format=json&pretty", esURL)
	if detailedFlag {
		url = fmt.Sprintf("%s&detailed", url)
	}
	t.SetCaption(url)

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

			row := table.Row{
				node.Name,
				fmt.Sprintf("%v:%v", task.Node, task.Id),
				task.Action,
				task.Type,
				task.StartTimeInMillis,
				task.RunningTimeInNanos,
				duration,
			}

			if detailedFlag {
				row = append(row, task.Description)
			}

			tableRows = append(tableRows, row)
		}
	}

	sort.Slice(tableRows, func(i, j int) bool {
		return tableRows[i][6].(time.Duration) < tableRows[j][6].(time.Duration)
	})

	t.AppendRows(tableRows)
	fmt.Println(render(t))
	return nil

}
