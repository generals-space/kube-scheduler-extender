package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/emicklei/go-restful"
	"k8s.io/klog/v2"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/apis/extender/v1"

	"github.com/generals-space/kube-scheduler-extender/pkg/schedule"
)

const (
	apiPrefix      = "/scheduler"
	predicatesPath = apiPrefix + "/predicates"
	prioritiesPath = apiPrefix + "/priorities"
	preemptionPath = apiPrefix + "/preemption"
	bindPath       = apiPrefix + "/bind"
)

// registPPROF 注册 pprof 页面
func registPPROF() {
	// 一个 WebService 表示一个前缀对象, 比如 /user, 通过 Path() 方法设置.
	// 之后可以通过 ws.Route() 在这个前缀下添加各种操作.
	pprofBasePath := "/debug"
	// 设置 pprof 的基本路径前缀
	ws := new(restful.WebService).Path(pprofBasePath)

	handlePprofEndpoint := func(req *restful.Request, resp *restful.Response) {
		name := strings.TrimPrefix(req.Request.URL.Path, pprofBasePath)
		switch name {
		case "profile":
			pprof.Profile(resp, req.Request)
		case "symbol":
			pprof.Symbol(resp, req.Request)
		case "cmdline":
			pprof.Cmdline(resp, req.Request)
		case "trace":
			pprof.Trace(resp, req.Request)
		default:
			pprof.Index(resp, req.Request)
		}
	}

	ws.Route(ws.GET("/{subpath:*}").To(func(req *restful.Request, resp *restful.Response) {
		handlePprofEndpoint(req, resp)
	})).Doc("pprof endpoint")

	restful.Add(ws)
}

// PredicateWrapper 将 Predicate 对象封装成 go-restful 可以使用的 handler 类型
func PredicateWrapper(predicate schedule.Predicate) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		var buf bytes.Buffer
		body := io.TeeReader(req.Request.Body, &buf)
		klog.Infof("predicate %s extender args: %s", predicate.Name, buf.String())

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
			extenderFilterResult = predicate.Handler(extenderArgs)
		}

		klog.Infof(
			"predicate %s extender filter result: %v",
			predicate.Name, extenderFilterResult.Nodes,
		)
		resp.WriteAsJson(extenderFilterResult)
	}
}

// BindWrapper ...
func BindWrapper() restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
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
}

// PreemptionWrapper ...
func PreemptionWrapper() restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
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
}

// PrioritiesWrapper ...
func PrioritiesWrapper() restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
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
}

func main() {
	registPPROF()

	var ws *restful.WebService

	ws = &restful.WebService{}
	ws.Path(predicatesPath)
	predicates := schedule.ListPredicates()
	for _, p := range predicates {
		klog.Infof("regist route predicate: %s", p.Name)
		path := p.Name
		ws.Route(
			ws.POST(path).To(PredicateWrapper(p)),
		)
	}
	restful.Add(ws)

	ws = &restful.WebService{}
	ws.Route(
		ws.POST(bindPath).To(BindWrapper()),
	)
	ws.Route(
		ws.POST(prioritiesPath).To(PrioritiesWrapper()),
	)
	ws.Route(
		ws.POST(preemptionPath).To(PreemptionWrapper()),
	)
	restful.Add(ws)

	http.ListenAndServe(":8080", nil)
}
