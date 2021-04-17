package schedule

import (
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

func PrioritiesHandler(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	klog.Infof(
		"Priorities request: Pod: %+v, Nodes: %+v, NodeNames: %+v",
		args.Pod, args.Nodes, args.NodeNames,
	)

	// 注意这里 priorityList 的声明和使用方式, 如果直接声明为 HostPriorityList{} 对象,
	// 是没有办法通过 append 为其追加元素的.
	var priorityList schedulerapi.HostPriorityList
	priorityList = make([]schedulerapi.HostPriority, len(args.Nodes.Items))

	for i, node := range args.Nodes.Items {
		priorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: 0,
		}
	}
	return &priorityList
}
