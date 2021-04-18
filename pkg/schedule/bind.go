package schedule

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"

	"github.com/generals-space/kube-scheduler-extender/pkg/k8s"
)

// Bind ...
func Bind(args *schedulerapi.ExtenderBindingArgs) (result *schedulerapi.ExtenderBindingResult) {
	klog.Infof(
		"binding request: podName: %s, namespace: %s, uid: %s, node: %s",
		args.PodName, args.PodNamespace, args.PodUID, args.Node,
	)
	var err error
	result = &schedulerapi.ExtenderBindingResult{}
	podName, ns := args.PodName, args.PodNamespace

	binding := &corev1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      podName,
			UID:       args.PodUID,
		},
		Target: corev1.ObjectReference{
			Kind: "Node",
			Name: args.Node,
		},
	}

	err = k8s.KubeClient.Client.CoreV1().Pods(ns).Bind(binding)
	if err != nil {
		klog.Errorf(
			"failed to bind %s/%s to %s: %s",
			ns, podName, args.Node, err,
		)
		result.Error = err.Error()
		return
	}

	pod, err := k8s.KubeClient.Client.CoreV1().Pods(ns).Get(podName, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("failed to get pod: %s/%s after binding: %s", ns, podName, err)
		result.Error = err.Error()
		return
	}

	klog.V(5).Infof("get pod after binding: %+v", pod)

	return
}
