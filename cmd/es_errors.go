package cmd

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
