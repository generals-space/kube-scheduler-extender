package schedule

import (
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"
)

func PreemptionHandler(args schedulerapi.ExtenderPreemptionArgs) *schedulerapi.ExtenderPreemptionResult {
	klog.Infof(
		"preemption request: Pod: %+v, NodeNameToVictims: %+v, NodeNameToMetaVictims: %+v",
		args.Pod, args.NodeNameToVictims, args.NodeNameToMetaVictims,
	)
	return &schedulerapi.ExtenderPreemptionResult{
		NodeNameToMetaVictims: args.NodeNameToMetaVictims,
	}
}
