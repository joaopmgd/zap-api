package main

import (
	"os"

	"gitlab.com/zap-api/app"
	"gitlab.com/zap-api/config"
)

func main() {
	config := config.GetConfig()
	app := &app.App{}
	app.Initialize(config)
	app.Run(os.Getenv("HOST"))
}
