package main

import (
	"os"

	"gitlab.com/api/app"
	"gitlab.com/api/config"
)

func main() {
	config := config.GetConfig()
	app := &app.App{}
	app.Initialize(config)
	app.Run(os.Getenv("HOST"))
}
