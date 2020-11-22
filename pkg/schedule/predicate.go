package schedule

import (
	apicorev1 "k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

// Predicate ...
type Predicate struct {
	Name string
	Func func(pod apicorev1.Pod, node apicorev1.Node) (bool, error)
}

var (
	alwaysTruePredicate = Predicate{
		Name: "always_true",
		Func: func(pod apicorev1.Pod, node apicorev1.Node) (bool, error) {
			return true, nil
		},
	}
)

// Handler ...
func (p Predicate) Handler(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	predicateNodes := args.Nodes.Items
	return &schedulerapi.ExtenderFilterResult{
		Nodes: &apicorev1.NodeList{
			Items: predicateNodes,
		},
		Error: "",
	}
}

// ListPredicates 返回自定义的 predicate 规则列表, 由主调函数注册到 http server
func ListPredicates() (predicates []Predicate) {
	predicates = append(predicates, alwaysTruePredicate)
	return
}
