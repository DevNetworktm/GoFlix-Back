package main

import (
	"fmt"
	"perso.go/GoFlix-Back/packages/easyapi"
	"perso.go/GoFlix-Back/packages/easypostgres"
	ApiRoutes "perso.go/GoFlix-Back/routes"
)

func main() {
	app := easyapi.New(8080)
	p, err := easypostgres.Open("dev", "1234", "goflix_project")
	if err != nil {
		app.HandlerError(app.Error.ERROR_FATAL, fmt.Sprintf("an error occurred open db, err=\"%v\"", err))
	}

	err = InitTables(p)
	if err != nil {
		app.HandlerError(app.Error.ERROR_FATAL, fmt.Sprintf("an error occurred init schema db, err=\"%v\"", err))
	} else {
		app.Logger.Log("All Schema is created !")
	}

	app.SetDatabase(p)

	app.JWT.New("secretKey", "secretRefreshKey")

	// Routes
	app.Router.AddRoute("/movies", ApiRoutes.MoviesRoutes(app))
	app.Router.AddRoute("/auth", ApiRoutes.AuthRoutes(app))

	app.Start()
}

func InitTables(p *easypostgres.PostgreSQL) error {
	const schemaMovie string = `
		CREATE TABLE IF NOT EXISTS movies (
			id SERIAL PRIMARY KEY,
			title TEXT,
			release_date DATE,
			duration INTEGER,
			synopsis TEXT,
			realisator TEXT[],
			productor TEXT[],
			actor TEXT[],
			picture TEXT,
			trailer_url TEXT
		);
	`

	const schemaAdmin string = `
		CREATE TABLE IF NOT EXISTS admin (
			id SERIAL PRIMARY KEY,
			username TEXT,
			email TEXT,
			password TEXT,
			avatar TEXT,
			created_at TEXT   
		)
	`

	err := easypostgres.NewExecInit(p, schemaMovie).Exec()
	if err != nil {
		return err
	}

	err = easypostgres.NewExecInit(p, schemaAdmin).Exec()
	if err != nil {
		return err
	}
	return nil
}
