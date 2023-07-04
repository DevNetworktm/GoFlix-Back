package ApiController

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"perso.go/GoFlix-Back/packages/easyapi"
	"perso.go/GoFlix-Back/packages/easyapi/controllers"
	"perso.go/GoFlix-Back/packages/easyapi/manager"
	"perso.go/GoFlix-Back/packages/easypostgres"
	"time"
)

func AuthSignUpController(app *easyapi.App) *controllers.ControllerManager {
	body := &easypostgres.Admin{}
	return controllers.New(func(request *manager.Request, response *manager.Response) {

		exist, err := app.Db.VerifyExist(body.Email)
		if err != nil {
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		} else if exist {
			response.SendStatus(http.StatusConflict)
			return
		}

		hasher := app.NewHasher(body.Password, "")
		newPassword, err := hasher.Hasher()
		if err != nil {
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		}

		body.Password = newPassword
		body.CreatedAt = time.DateTime

		err = app.Db.SignUp(body)
		if err != nil {
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		}

		response.Status(http.StatusCreated).Send("you were well and truly recorded !")
	}, body)
}

func AuthLoginController(app *easyapi.App) *controllers.ControllerManager {
	body := &easypostgres.Admin{}
	return controllers.New(func(request *manager.Request, response *manager.Response) {
		exist, err := app.Db.VerifyExist(body.Email)
		if err != nil {
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		} else if !exist {
			response.SendStatus(http.StatusNotFound)
			return
		}

		password := body.Password

		admin, err := app.Db.Login(body)
		if err != nil {
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		}

		hasher := app.NewHasher(password, admin.Password)
		if !hasher.Verify() {
			response.SendStatus(http.StatusNotFound)
			return
		}

		adminJ := admin.FormatAdminToJWT(time.Now().Add((time.Hour * 60 * 3) * time.Duration(1)))

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, adminJ)
		ss, err := token.SignedString([]byte(app.JWT.GetPrivateKey()))
		if err != nil {
			response.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return
		}

		res := map[string]string{
			"token": ss,
		}

		response.Status(http.StatusOK).Json(res)
	}, body)
}
