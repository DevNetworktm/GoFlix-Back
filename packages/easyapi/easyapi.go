package easyapi

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"perso.go/GoFlix-Back/packages/easyapi/router"
	"perso.go/GoFlix-Back/packages/easypostgres"
)

// EasyApi Func
type App struct {
	Port   int16
	Error  Error
	Logger Logger

	Router *router.Router
	Db     *easypostgres.PostgreSQL
	JWT    *JWT
}

func New(port int16) *App {
	m := make(chan string)

	app := &App{
		Port: port,
		Error: Error{
			channel:       m,
			ERROR_FATAL:   0,
			ERROR_WARNING: 0,
		},
		Logger: Logger{
			channel: m,
		},
		Router: router.New(),
		JWT: &JWT{
			privateKey:        "",
			privateRefreshKey: "",
		},
	}

	go writeMessage(m)

	return app
}

func (app *App) SetDatabase(bdd *easypostgres.PostgreSQL) {
	app.Db = bdd
}

func (app *App) Start() {
	app.Logger.Log(fmt.Sprintf("Easy Api was start in http://192.168.1.153:%d", app.Port))

	app.Router.ListenRouter(app.Logger.channel)

	http.HandleFunc("/", app.Router.SearchRoutes(app.Logger.channel))

	err := http.ListenAndServe(fmt.Sprintf("192.168.1.153:%d", app.Port), nil)
	if err == nil {
		app.HandlerError(app.Error.ERROR_FATAL, fmt.Sprintf("an error occurred during api launch, err=\"%v\"", err))
	}
}

func (app *App) HandlerError(errorType int8, message string) {
	switch errorType {
	case app.Error.ERROR_FATAL:
		app.Error.ErrorF(message)
	case app.Error.ERROR_WARNING:
		app.Error.ErrorW(message)
	}
}

// Error
type Error struct {
	channel       chan string
	ERROR_FATAL   int8
	ERROR_WARNING int8
}

func (e *Error) ErrorF(message string) {
	e.channel <- fmt.Sprintf("[FATAL] %s", message)
}

func (e *Error) ErrorW(message string) {
	e.channel <- fmt.Sprintf("[FATAL] %s", message)
}

// Logger
type Logger struct {
	channel chan string
}

func (l *Logger) Log(message string) {
	l.channel <- fmt.Sprintf("[INFO] %s", message)
}

// ChanFunc
func writeMessage(c chan string) {
	for {
		log.Println(<-c)
	}
}

// Hasher
type Hasher struct {
	Password string
	Hash     string
	salt     int
}

func (app App) NewHasher(password string, hash string) *Hasher {
	h := &Hasher{
		salt:     10,
		Password: password,
		Hash:     hash,
	}

	return h
}

func (h *Hasher) Hasher() (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(h.Password), h.salt)
	return string(bytes), err
}

func (h *Hasher) Verify() bool {
	err := bcrypt.CompareHashAndPassword([]byte(h.Hash), []byte(h.Password))
	return err == nil
}

// JWT

type JWT struct {
	privateKey        string
	privateRefreshKey string
}

func (j *JWT) New(privateKey, privateRefreshKey string) {
	j.privateKey = privateKey
	j.privateRefreshKey = privateRefreshKey
}

func (j *JWT) GetPrivateKey() string {
	return j.privateKey
}

func (j *JWT) GetPrivateRefreshKey() string {
	return j.privateRefreshKey
}
