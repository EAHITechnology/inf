package enet

import (
	"net/http"

	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	//路由器
	Router *httprouter.Router
	//pprof 开启开关
	PProf bool
}

func NewRouter() (*Router, error) {
	r := httprouter.New()
	return &Router{Router: r}, nil
}

/*
打开pprof开关
*/
func (this *Router) OpenPProf() {
	this.PProf = true
}

/*
get路由
*/
func (this *Router) Get(path string, h http.Handler) {
	this.Router.Handler(http.MethodGet, path, h)
}

/*
post路由
*/
func (this *Router) Post(path string, h http.Handler) {
	this.Router.Handler(http.MethodPost, path, h)
}

/*
get方法路由
*/
func (this *Router) GetFunc(path string, handler http.HandlerFunc) {
	this.Router.HandlerFunc(http.MethodGet, path, handler)
}

/*
post方法路由
*/
func (this *Router) PostFunc(path string, handler http.HandlerFunc) {
	this.Router.Handler(http.MethodPost, path, handler)
}

func SetHttpLisPort(httpPort string) {
	HTTP_PORT = httpPort
}

/*
异步监听，程序不会阻塞在这里
*/
func (this *Router) Listen() {
	if this.PProf {
		this.Router.HandlerFunc(http.MethodGet, "/debug/pprof", pprof.Index)
		this.Router.HandlerFunc(http.MethodGet, "/debug/cmdline", pprof.Cmdline)
		this.Router.HandlerFunc(http.MethodGet, "/debug/profile", pprof.Profile)
		this.Router.HandlerFunc(http.MethodGet, "/debug/trace", pprof.Trace)
		this.Router.Handler(http.MethodGet, "/debug/block", pprof.Handler("block"))
		this.Router.Handler(http.MethodGet, "/debug/allocs", pprof.Handler("allocs"))
		this.Router.Handler(http.MethodGet, "/debug/goroutine", pprof.Handler("goroutine"))
		this.Router.Handler(http.MethodGet, "/debug/heap", pprof.Handler("heap"))
		this.Router.Handler(http.MethodGet, "/debug/mutex", pprof.Handler("mutex"))
		this.Router.Handler(http.MethodGet, "/debug/threadcreate", pprof.Handler("threadcreate"))
	}
	go func() {
		if err := http.ListenAndServe(":"+HTTP_PORT, this.Router); err != nil {
			panic("Listen ListenAndServe err:" + err.Error())
		}
	}()
}
