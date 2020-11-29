package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var allocationExplainCmd = &cobra.Command{
	Use:   "allocation",
	Short: "allocation",
	Long:  "allocation",
	Run:   clusterAllocatoinExplain,
}

const (
	allocPath = "/_cluster/allocation/explain"
)

// {
//     "error": {
//         "root_cause": [
//             {
//                 "type": "illegal_argument_exception",
//                 "reason": "unable to find any unassigned shards to explain [ClusterAllocationExplainRequest[useAnyUnassignedShard=true,includeYesDecisions?=false]"
//             }
//         ],
//         "type": "illegal_argument_exception",
//         "reason": "unable to find any unassigned shards to explain [ClusterAllocationExplainRequest[useAnyUnassignedShard=true,includeYesDecisions?=false]"
//     },
//     "status": 400
// }

type esRootError struct {
	Errors     esError `json:"error"`
	StatusCode int     `json:"status"`
}

type esError struct {
	Causes    []cause `json:"root_cause"`
	ErrorType string  `json:"type"`
	Reason    string  `json:"reason"`
}

type cause struct {
	ErrorType string `json:"type"`
	Reason    string `json:"reason"`
}

func init() {
	clusterCmd.AddCommand(allocationExplainCmd)
}

func clusterAllocatoinExplain(cmd *cobra.Command, args []string) {

	// ctx := context.Background()

	u, err := url.Parse(esURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	u.Path = path.Join(u.Path, allocPath)

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode == http.StatusBadRequest {
		var foo esRootError
		err = json.Unmarshal(body, &foo)
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.HasPrefix(foo.Errors.Reason, "unable to find any unassigned shards to explain") {
			fmt.Println("no unassigned shards")
			return
		}

		for _, r := range foo.Errors.Causes {
			fmt.Println("reson", r.ErrorType, r.Reason)
		}

		return
	}

	fmt.Println(string(body))
	return
}
