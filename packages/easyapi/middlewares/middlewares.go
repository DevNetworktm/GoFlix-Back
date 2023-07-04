package middlewares

import "perso.go/GoFlix-Back/packages/easyapi/manager"

type Exec func(req *manager.Request, res *manager.Response) (next bool, finish bool, err error)

type MiddlewaresManager struct {
	exec Exec
	body interface{}
}

func New(exec Exec, body interface{}) *MiddlewaresManager {
	c := &MiddlewaresManager{
		exec: exec,
		body: body,
	}

	return c
}

func (c *MiddlewaresManager) GetExec() Exec {
	return c.exec
}

func (c *MiddlewaresManager) GetBody() interface{} {
	return c.body
}
