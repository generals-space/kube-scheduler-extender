package schedule

import (
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

func BindHandler(args schedulerapi.ExtenderBindingArgs) *schedulerapi.ExtenderBindingResult {
	klog.Infof(
		"binding request: podName: %s, namespace: %s, uid: %s, node: %s",
		args.PodName, args.PodNamespace, args.PodUID, args.Node,
	)
	return &schedulerapi.ExtenderBindingResult{
		Error: "",
	}
}
