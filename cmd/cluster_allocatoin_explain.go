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
	Use:           "allocation",
	Short:         "allocation",
	Long:          "allocation",
	RunE:          clusterAllocatoinExplain,
	SilenceErrors: true,
	SilenceUsage:  true,
}

const (
	allocPath = "/_cluster/allocation/explain"
)

func init() {
	clusterCmd.AddCommand(allocationExplainCmd)
}

func clusterAllocatoinExplain(cmd *cobra.Command, args []string) error {

	u, err := url.Parse(esURL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, allocPath)

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		var foo esRootError
		err = json.Unmarshal(body, &foo)
		if err != nil {
			return err
		}

		if strings.HasPrefix(foo.Errors.Reason, "unable to find any unassigned shards to explain") {
			fmt.Println("no unassigned shards")
			return nil
		}

		for _, r := range foo.Errors.Causes {
			fmt.Println("reson", r.ErrorType, r.Reason)
		}

		return nil
	}

	fmt.Println(string(body))
	return nil
}
