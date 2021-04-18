package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/emicklei/go-restful"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"

	"github.com/generals-space/kube-scheduler-extender/pkg/route"
	"github.com/generals-space/kube-scheduler-extender/pkg/schedule"
)

const (
	apiPrefix      = "/scheduler"
	predicatesPath = "/predicates"
	prioritiesPath = "/priorities"
	preemptionPath = "/preemption"
	bindPath       = "/bind"
)

// PredicateHandler ...
func PredicateHandler(req *restful.Request, resp *restful.Response) {
	var buf bytes.Buffer
	body := io.TeeReader(req.Request.Body, &buf)

	var extenderArgs schedulerapi.ExtenderArgs
	var extenderFilterResult *schedulerapi.ExtenderFilterResult

	err := json.NewDecoder(body).Decode(&extenderArgs)
	if err != nil {
		extenderFilterResult = &schedulerapi.ExtenderFilterResult{
			Nodes:       nil,
			FailedNodes: nil,
			Error:       err.Error(),
		}
	} else {
		extenderFilterResult = schedule.PredicateHandler(extenderArgs)
	}
	resp.WriteAsJson(extenderFilterResult)
}

func BindHandler(req *restful.Request, resp *restful.Response) {
	var buf bytes.Buffer
	body := io.TeeReader(req.Request.Body, &buf)

	var extenderBindingArgs schedulerapi.ExtenderBindingArgs
	var extenderBindingResult *schedulerapi.ExtenderBindingResult

	err := json.NewDecoder(body).Decode(&extenderBindingArgs)
	if err != nil {
		extenderBindingResult = &schedulerapi.ExtenderBindingResult{
			Error: err.Error(),
		}
	} else {
		extenderBindingResult = schedule.BindHandler(extenderBindingArgs)
	}
	resp.WriteAsJson(extenderBindingResult)
}

func PreemptionHandler(req *restful.Request, resp *restful.Response) {
	var buf bytes.Buffer
	body := io.TeeReader(req.Request.Body, &buf)

	var extenderPreemptionArgs schedulerapi.ExtenderPreemptionArgs
	var extenderPreemptionResult *schedulerapi.ExtenderPreemptionResult

	err := json.NewDecoder(body).Decode(&extenderPreemptionArgs)
	if err != nil {
		extenderPreemptionResult = &schedulerapi.ExtenderPreemptionResult{}
	} else {
		extenderPreemptionResult = schedule.PreemptionHandler(extenderPreemptionArgs)
	}
	resp.WriteAsJson(extenderPreemptionResult)
}

func PrioritiesHandler(req *restful.Request, resp *restful.Response) {
	var buf bytes.Buffer
	body := io.TeeReader(req.Request.Body, &buf)

	var extenderArgs schedulerapi.ExtenderArgs
	var extenderPreemptionResult *schedulerapi.HostPriorityList

	err := json.NewDecoder(body).Decode(&extenderArgs)
	if err != nil {
		extenderPreemptionResult = &schedulerapi.HostPriorityList{}
	} else {
		extenderPreemptionResult = schedule.PrioritiesHandler(extenderArgs)
	}
	resp.WriteAsJson(extenderPreemptionResult)
}

func main() {
	route.RegistPPROF()

	ws := &restful.WebService{}
	ws.Path(apiPrefix)
	ws.Route(ws.POST(predicatesPath).To(PredicateHandler))
	ws.Route(ws.POST(prioritiesPath).To(PrioritiesHandler))
	ws.Route(ws.POST(bindPath).To(BindHandler))
	ws.Route(ws.POST(preemptionPath).To(PreemptionHandler))
	restful.Add(ws)

	http.ListenAndServe(":8080", nil)
}
