package main

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"k8s.io/klog/v2"

	"github.com/generals-space/kube-scheduler-extender/pkg/k8s"
	"github.com/generals-space/kube-scheduler-extender/pkg/route"
)

const (
	apiPrefix      = "/scheduler"
	predicatesPath = "/predicates"
	prioritizePath = "/prioritize"
	preemptionPath = "/preemption"
	bindPath       = "/bind"
)

func main() {
	klog.Infof("init kube client")
	k8s.InitKubeClient()

	klog.Infof("regist pprof handler")
	route.RegistPPROF()

	klog.Infof("regist schedule extender handler")
	ws := &restful.WebService{}
	ws.Path(apiPrefix)
	// 预选过滤接口
	ws.Route(ws.POST(predicatesPath).To(route.PredicateHandler))
	// 优选打分接口
	ws.Route(ws.POST(prioritizePath).To(route.PrioritizeHandler))
	ws.Route(ws.POST(preemptionPath).To(route.PreemptionHandler))
	// 当核心 scheduler 调度器确定 Node 与 Pod 时的回调接口,
	// 核心调度器会把即将绑定的一对 Pod 与 Node 发送到这个接口.
	ws.Route(ws.POST(bindPath).To(route.BindHandler))
	restful.Add(ws)

	klog.Infof("start listening")
	http.ListenAndServe(":8080", nil)
}
