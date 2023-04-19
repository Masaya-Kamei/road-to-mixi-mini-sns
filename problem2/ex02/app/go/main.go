package main

import (
	"problem1/configs"
	"problem1/controllers"
	"problem1/models"
	"strconv"
)

func main() {
	models.InitDb()
	defer models.CloseDb()

	models.InitCache()

	e := controllers.NewRouter()
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(configs.Get().Server.Port)))
}
