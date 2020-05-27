package main

import (
	"github.com/paulosabaini/comidamexicana/app"
	"github.com/paulosabaini/comidamexicana/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
