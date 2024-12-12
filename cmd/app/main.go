package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/xoticdsign/returnauf/config"
	_ "github.com/xoticdsign/returnauf/docs"
	"github.com/xoticdsign/returnauf/internal/app"
)

// Общее описание
//
// @title                      returnauf
// @version                    1.0.0
// @description                TODO
// @contact.name               xoti$
// @contact.url                https://t.me/xoticdsign
// @contact.email              xoticdollarsign@outlook.com
// @license.name               MIT
// @license.url                https://mit-license.org/
// @host                       127.0.0.1:8080
// @BasePath                   /
// @produce                    json
// @schemes                    http
//
// @securitydefinitions.apikey KeyAuth
// @in                         query
// @name                       returnauf-key
func main() {
	godotenv.Load()

	conf := config.LoadConfig()

	app, err := app.InitApp(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Listen(conf.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}
}
