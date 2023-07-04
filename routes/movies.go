package ApiRoutes

import (
	ApiController "perso.go/GoFlix-Back/controllers"
	ApiMiddlewares "perso.go/GoFlix-Back/middlewares"
	"perso.go/GoFlix-Back/packages/easyapi"
	"perso.go/GoFlix-Back/packages/easyapi/router"
)

func MoviesRoutes(app *easyapi.App) *router.ChildrenRouter {
	r := router.NewChildrenRouter()

	r.Get("", ApiController.GetAllMovieController(app))
	r.Get("/:id", ApiController.GetOneMovieController(app))

	r.Post("", ApiController.CreateMovieController(app), ApiMiddlewares.JWTMiddleware(app))

	return r
}
