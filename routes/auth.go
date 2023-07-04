package ApiRoutes

import (
	ApiController "perso.go/GoFlix-Back/controllers"
	ApiMiddlewares "perso.go/GoFlix-Back/middlewares"
	"perso.go/GoFlix-Back/packages/easyapi"
	"perso.go/GoFlix-Back/packages/easyapi/router"
)

func AuthRoutes(app *easyapi.App) *router.ChildrenRouter {
	r := router.NewChildrenRouter()

	r.Post("/signup", ApiController.AuthSignUpController(app), ApiMiddlewares.JWTMiddleware(app))
	r.Post("/login", ApiController.AuthLoginController(app))

	return r
}
