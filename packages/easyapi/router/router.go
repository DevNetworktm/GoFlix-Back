package router

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"perso.go/GoFlix-Back/packages/easyapi/controllers"
	"perso.go/GoFlix-Back/packages/easyapi/manager"
	"perso.go/GoFlix-Back/packages/easyapi/middlewares"
	"strings"
)

type Router struct {
	Methods []*Method
}

type Request struct {
}

type ChildrenRouter struct {
	Methods []*Method
}

func New() *Router {
	r := &Router{
		Methods: []*Method{},
	}

	return r
}

func NewChildrenRouter() *ChildrenRouter {
	rc := &ChildrenRouter{}

	return rc
}

func (router *Router) AddRoute(path string, routerChildrenFunc *ChildrenRouter) {
	rc := routerChildrenFunc

	for _, m := range rc.Methods {
		var newPath string

		if path == "/" {
			newPath = m.path
		} else {
			newPath = path + m.path
		}

		router.Methods = append(router.Methods, &Method{
			path:        newPath,
			method:      m.method,
			manager:     m.manager,
			middlewares: m.middlewares,
		})
	}
}

func (router *Router) ListenRouter(c chan string) {

	c <- "\n\n"
	c <- "[ROUTER] List Of Request :"
	for _, method := range router.Methods {
		m := ""
		p := ""

		switch method.method {
		case "GET":
			m = color.New(color.FgBlue).SprintfFunc()("GET")
			p = color.New(color.FgBlue).SprintfFunc()(method.path)
		case "POST":
			m = color.New(color.FgGreen).SprintfFunc()("POST")
			p = color.New(color.FgGreen).SprintfFunc()(method.path)
		case "PUT":
			m = color.New(color.FgYellow).SprintfFunc()("PUT")
			p = color.New(color.FgYellow).SprintfFunc()(method.path)
		case "DELETE":
			m = color.New(color.FgRed).SprintfFunc()("DELETE")
			p = color.New(color.FgRed).SprintfFunc()(method.path)
		case "PATCH":
			m = color.New(color.FgHiMagenta).SprintfFunc()("PATCH")
			p = color.New(color.FgHiMagenta).SprintfFunc()(method.path)
		}

		c <- fmt.Sprintf("[ROUTER] [%s] => %s is initialised", m, p)
	}
}

func (router *Router) SearchRoutes(c chan string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(200)
			return
		}
		url := strings.Split(r.RequestURI, "?")[0]
		urlSplit := strings.Split(url, "/")[1:]

		req := &manager.Request{
			Header:        r.Header,
			ContentLength: r.ContentLength,
			Form:          r.Form,
			Method:        r.Method,
			Host:          r.Host,

			RequestHTTP: r,
		}

		res := &manager.Response{
			Write: w,
		}

		find := false

		for _, m := range router.Methods {
			query := r.URL.Query()
			params := map[string]string{}

			req.Query = query

			methodUrlSplit := strings.Split(m.path, "/")[1:]

			if r.Method == m.method && url == m.path {
				find = true
				if m.manager.GetBody() != nil {
					err := json.NewDecoder(r.Body).Decode(m.manager.GetBody())
					if err != nil {
						res.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
						return
					}
				}

				if len(m.middlewares) > 0 {
					for _, mi := range m.middlewares {
						next, finish, err := mi.GetExec()(req, res)
						if err != nil {
							res.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
							return
						}

						if finish {
							return
						}

						if next {
							continue
						}
					}
				}
				m.manager.GetExec()(req, res)
			} else if len(urlSplit) == len(methodUrlSplit) {
				err := false

				for index, part := range urlSplit {
					if err {
						continue
					}

					param := methodUrlSplit[index]

					if strings.Contains(methodUrlSplit[index], ":") {
						paramName := strings.ReplaceAll(param, ":", "")
						params[paramName] = part
						methodUrlSplit[index] = part
					} else if param == part {
						continue
					} else {
						err = true
					}
				}

				if err {
					continue
				}

				newUrl := "/" + strings.Join(methodUrlSplit, "/")
				if newUrl == url && r.Method == m.method {

					req.Params = params
					find = true
					if m.manager.GetBody() != nil {

						err := json.NewDecoder(r.Body).Decode(m.manager.GetBody())
						if err != nil {
							res.SendStatus(http.StatusInternalServerError)
							return
						}
					}

					if len(m.middlewares) > 0 {
						for _, mi := range m.middlewares {
							next, finish, err := mi.GetExec()(req, res)
							if err != nil {
								res.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
								return
							}

							if finish {
								return
							}

							if next {
								continue
							}
						}
					}
					m.manager.GetExec()(req, res)
					return
				}
			}
			continue
		}

		if !find {
			res.SendStatus(http.StatusNotFound)
		}
	}
}

// Methods
type Method struct {
	path        string
	manager     *controllers.ControllerManager
	method      string
	middlewares []*middlewares.MiddlewaresManager
}

func (rc *ChildrenRouter) Get(path string, manager *controllers.ControllerManager, middlewares ...*middlewares.MiddlewaresManager) {
	m := &Method{
		method:      "GET",
		manager:     manager,
		path:        path,
		middlewares: middlewares,
	}

	rc.Methods = append(rc.Methods, m)
}

func (rc *ChildrenRouter) Post(path string, manager *controllers.ControllerManager, middlewares ...*middlewares.MiddlewaresManager) {
	m := &Method{
		method:      "POST",
		manager:     manager,
		path:        path,
		middlewares: middlewares,
	}

	rc.Methods = append(rc.Methods, m)
}

func (rc *ChildrenRouter) Put(path string, manager *controllers.ControllerManager, middlewares ...*middlewares.MiddlewaresManager) {
	m := &Method{
		method:      "PUT",
		manager:     manager,
		path:        path,
		middlewares: middlewares,
	}

	rc.Methods = append(rc.Methods, m)
}

func (rc *ChildrenRouter) Delete(path string, manager *controllers.ControllerManager, middlewares ...*middlewares.MiddlewaresManager) {
	m := &Method{
		method:      "DELETE",
		manager:     manager,
		path:        path,
		middlewares: middlewares,
	}

	rc.Methods = append(rc.Methods, m)
}

func (rc *ChildrenRouter) Patch(path string, manager *controllers.ControllerManager, middlewares ...*middlewares.MiddlewaresManager) {
	m := &Method{
		method:      "PATCH",
		manager:     manager,
		path:        path,
		middlewares: middlewares,
	}

	rc.Methods = append(rc.Methods, m)
}

func (rc *ChildrenRouter) GetMethods() []*Method {
	return rc.Methods
}
