package controllers

import "perso.go/GoFlix-Back/packages/easyapi/manager"

type ControllerManager struct {
	exec func(req *manager.Request, res *manager.Response)
	body interface{}
}

func New(exec func(req *manager.Request, res *manager.Response), body interface{}) *ControllerManager {
	c := &ControllerManager{
		exec: exec,
		body: body,
	}

	return c
}

func (c *ControllerManager) GetExec() func(req *manager.Request, res *manager.Response) {
	return c.exec
}

func (c *ControllerManager) GetBody() interface{} {
	return c.body
}
