package schedule

import (
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

// ExtenderArgs.Pod: 待调度的 Pod 对象
// ExtenderArgs.Nodes: 可供调度的 Node 列表对象
// ExtenderArgs.NodeNames: 可供调度的 Node 名称列表
// Nodes 与 NodeNames 貌似不会同时出现, 通过 ExtenderConfig.NodeCacheCapable 进行配置, 
// 原理应该是, 如 extender 服务没有开启 LocalStorage 时, 核心 scheduler 的请求参数里传入的都是 Node 对象;
// 而当 extender 开启了 LocalStorage 后, 核心 scheduler 的请求就可以直接传入 Node 名称了.
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
