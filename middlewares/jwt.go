package ApiMiddlewares

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"perso.go/GoFlix-Back/packages/easyapi"
	"perso.go/GoFlix-Back/packages/easyapi/manager"
	"perso.go/GoFlix-Back/packages/easyapi/middlewares"
	"perso.go/GoFlix-Back/packages/easypostgres"
	"strings"
)

func JWTMiddleware(app *easyapi.App) *middlewares.MiddlewaresManager {
	return middlewares.New(func(req *manager.Request, res *manager.Response) (next bool, finish bool, err error) {
		Authorization := req.GetHeader("Authorization")

		if Authorization == "" {
			res.Status(http.StatusBadRequest).Send("Header: 'Authorization is missing !'")
			return false, true, nil
		} else if !strings.Contains(Authorization, "Bearer") {
			res.Status(http.StatusBadRequest).Send("Authorization does not belong to the format (\"Authorization\": \"Bearer [TOKEN]\")")
			return false, true, nil
		}

		if app.JWT.GetPrivateKey() == "" {
			res.SendStatus(http.StatusInternalServerError)
			app.HandlerError(app.Error.ERROR_FATAL, "If you want use JWT in your, you need to define 'app.JWT.New(privateKey string, privateRefreshKey string)'")
			return false, true, nil
		}

		Authorization = strings.Split(Authorization, "Bearer ")[1]

		token, err := jwt.ParseWithClaims(Authorization, &easypostgres.AdminJwt{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.JWT.GetPrivateKey()), nil
		})

		if err != nil {
			res.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return false, true, nil
		}

		if claims, ok := token.Claims.(*easypostgres.AdminJwt); ok && token.Valid {
			reqVar := map[string]interface{}{
				"id": claims.ID,
			}
			req.SetRequestVar(reqVar)
			return true, false, nil
		} else {
			res.Status(http.StatusInternalServerError).Send(fmt.Sprintf("%v", err))
			return false, true, nil
		}
	}, nil)
}
