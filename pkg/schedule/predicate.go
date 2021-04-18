package schedule

import (
	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/klog"

	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

// Predicate 对传入的 Node 对象列表进行过滤.
//
//	Parameters
//	- `ExtenderArgs.Pod`: 待调度的 Pod 对象
//	- `ExtenderArgs.Nodes`: 可供调度的 Node 列表对象
//	- `ExtenderArgs.NodeNames`: 可供调度的 Node 名称列表
//
//	Note
// Nodes 与 NodeNames 貌似不会同时出现, 通过 ExtenderConfig.NodeCacheCapable 进行配置.
// 原理应该是, 如 extender 服务没有开启 LocalStorage 时, 核心 scheduler 的请求参数里传入的都是 Node 对象;
// 而当 extender 开启了 LocalStorage 后, 核心 scheduler 的请求就可以直接传入 Node 名称了.
func Predicate(args *schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	nodeNamesInList := []string{}
	for _, node := range args.Nodes.Items {
		nodeNamesInList = append(nodeNamesInList, node.Name)
	}
	klog.Infof(
		"predicate request: Pod: %s, Nodes: %+v, NodeNames: %+v",
		args.Pod.Name, nodeNamesInList, args.NodeNames,
	)
	klog.V(5).Infof("predicate request: Pod: %+v, NodeList: %+v", args.Pod, args.Nodes)

	return &schedulerapi.ExtenderFilterResult{
		Nodes: &apicorev1.NodeList{
			Items: args.Nodes.Items,
		},
		Error: "",
	}
}
