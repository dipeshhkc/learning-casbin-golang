package main

import (
	"casbin-golang/model"
	"casbin-golang/route"
)

func main() {

	db, _ := model.DBConnection()
	route.SetupRoutes(db)
}
