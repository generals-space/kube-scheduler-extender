package route

import (
	"net/http/pprof"
	"strings"

	"github.com/emicklei/go-restful"
)

// RegistPPROF 注册 pprof 页面
func RegistPPROF() {
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
