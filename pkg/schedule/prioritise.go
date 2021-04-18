package schedule

import (
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

// Prioritise 为传入的 Node 对象列表进行打分.
func Prioritise(args *schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	nodeNamesInList := []string{}
	for _, node := range args.Nodes.Items {
		nodeNamesInList = append(nodeNamesInList, node.Name)
	}
	klog.Infof(
		"prioritize request: Pod: %s, Nodes: %+v, NodeNames: %+v",
		args.Pod.Name, nodeNamesInList, args.NodeNames,
	)
	klog.V(5).Infof("prioritize request: Pod: %+v, NodeList: %+v", args.Pod, args.Nodes)

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
