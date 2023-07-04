package ApiController

import (
	"fmt"
	"log"
	"net/http"
	"perso.go/GoFlix-Back/packages/easyapi"
	"perso.go/GoFlix-Back/packages/easyapi/controllers"
	"perso.go/GoFlix-Back/packages/easyapi/manager"
	"perso.go/GoFlix-Back/packages/easypostgres"
	"strconv"
)

func GetAllMovieController(app *easyapi.App) *controllers.ControllerManager {
	return controllers.New(func(request *manager.Request, response *manager.Response) {
		movies, err := app.Db.GetMovies()
		if err != nil {
			app.HandlerError(app.Error.ERROR_WARNING, fmt.Sprintf("%v", err))
		}

		response.Status(200).Json(movies)
	}, nil)
}

func GetOneMovieController(app *easyapi.App) *controllers.ControllerManager {
	return controllers.New(func(request *manager.Request, response *manager.Response) {
		id := request.Params["id"]
		num, err := strconv.Atoi(id)
		if err != nil {
			response.SendStatus(http.StatusBadRequest)
			return
		}

		movie, err := app.Db.GetMovie(num)
		if err != nil {
			app.HandlerError(app.Error.ERROR_WARNING, fmt.Sprintf("%v", err))
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		}

		response.Status(200).Json(movie)
	}, nil)
}

func CreateMovieController(app *easyapi.App) *controllers.ControllerManager {
	body := easypostgres.Movie{}
	return controllers.New(func(request *manager.Request, response *manager.Response) {
		id, err := app.Db.CreateMovie(&body)
		log.Println("tes")

		if err != nil {
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		}

		movie, err := app.Db.GetMovie(int(id))
		if err != nil {
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		}

		response.Status(http.StatusCreated).Json(movie)
	}, &body)
}
