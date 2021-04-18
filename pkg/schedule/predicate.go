package schedule

import (
	apicorev1 "k8s.io/api/core/v1"

	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

// PredicateHandler ...
func PredicateHandler(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	predicateNodes := args.Nodes.Items
	return &schedulerapi.ExtenderFilterResult{
		Nodes: &apicorev1.NodeList{
			Items: predicateNodes,
		},
		Error: "",
	}
}
