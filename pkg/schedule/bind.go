package schedule

import (
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

// BindHandler 算是核心 scheduler 调度器确定 Node 与 Pod 后的回调接口, 
// 核心调度器会把即将绑定的一对 Pod 与 Node 发送到这个接口.
func BindHandler(args schedulerapi.ExtenderBindingArgs) *schedulerapi.ExtenderBindingResult {
	klog.Infof(
		"binding request: podName: %s, namespace: %s, uid: %s, node: %s",
		args.PodName, args.PodNamespace, args.PodUID, args.Node,
	)
	return &schedulerapi.ExtenderBindingResult{
		Error: "",
	}
}
