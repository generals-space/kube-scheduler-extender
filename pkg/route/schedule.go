package route

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/emicklei/go-restful"
	"k8s.io/klog/v2"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"

	"github.com/generals-space/kube-scheduler-extender/pkg/schedule"
)

// PredicateHandler 解析来自核心 scheduler 的 POST 请求体, 并调用我们实际的预选方法进行筛选.
func PredicateHandler(req *restful.Request, resp *restful.Response) {
	klog.Infof("predicate request...")

	var buf bytes.Buffer
	body := io.TeeReader(req.Request.Body, &buf)

	args := &schedulerapi.ExtenderArgs{}
	var filterResult *schedulerapi.ExtenderFilterResult

	err := json.NewDecoder(body).Decode(args)
	if err != nil {
		filterResult = &schedulerapi.ExtenderFilterResult{
			Nodes:       nil,
			FailedNodes: nil,
			Error:       err.Error(),
		}
	} else {
		filterResult = schedule.Predicate(args)
	}
	klog.Infof("predicate response: %+v", filterResult)
	resp.WriteAsJson(filterResult)
}

// PrioritizeHandler 解析来自核心 scheduler 的 POST 请求体, 并调用我们实际的优化方法进行打分.
func PrioritizeHandler(req *restful.Request, resp *restful.Response) {
	klog.Infof("prioritie request...")

	var buf bytes.Buffer
	body := io.TeeReader(req.Request.Body, &buf)

	args := &schedulerapi.ExtenderArgs{}
	var hostPriorityList *schedulerapi.HostPriorityList

	err := json.NewDecoder(body).Decode(args)
	if err != nil {
		hostPriorityList = &schedulerapi.HostPriorityList{}
	} else {
		hostPriorityList = schedule.Prioritise(args)
	}
	klog.Infof("prioritie response: %+v", hostPriorityList)
	resp.WriteAsJson(hostPriorityList)
}

// PreemptionHandler ...
func PreemptionHandler(req *restful.Request, resp *restful.Response) {
	klog.Infof("preempt request...")

	var buf bytes.Buffer
	body := io.TeeReader(req.Request.Body, &buf)

	preemptionArgs := &schedulerapi.ExtenderPreemptionArgs{}
	var preemptionResult *schedulerapi.ExtenderPreemptionResult

	err := json.NewDecoder(body).Decode(preemptionArgs)
	if err != nil {
		preemptionResult = &schedulerapi.ExtenderPreemptionResult{}
	} else {
		preemptionResult = schedule.Preempt(preemptionArgs)
	}
	klog.Infof("preempt response: %+v", preemptionResult)
	resp.WriteAsJson(preemptionResult)
}

func BindHandler(req *restful.Request, resp *restful.Response) {
	klog.Infof("bind request...")

	var buf bytes.Buffer
	body := io.TeeReader(req.Request.Body, &buf)

	var bindingResult *schedulerapi.ExtenderBindingResult
	bindingArgs := &schedulerapi.ExtenderBindingArgs{}

	err := json.NewDecoder(body).Decode(bindingArgs)
	if err != nil {
		bindingResult = &schedulerapi.ExtenderBindingResult{
			Error: err.Error(),
		}
	} else {
		bindingResult = schedule.Bind(bindingArgs)
	}
	klog.Infof("bind response: %+v", bindingResult)
	resp.WriteAsJson(bindingResult)
}
