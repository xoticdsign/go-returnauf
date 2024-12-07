package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/xoticdsign/auf-citaty/config"
	_ "github.com/xoticdsign/auf-citaty/docs"
	"github.com/xoticdsign/auf-citaty/internal/app"
)

// Общее описание
//
// @title                      Auf Citaty API
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
// @name                       auf-citaty-key
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
