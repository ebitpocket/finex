package main

import (
	"github.com/nusa-exchange/finex/config"
	"github.com/nusa-exchange/finex/routes"
)

func main() {
	if err := config.InitializeConfig(); err != nil {
		config.Logger.Error(err.Error())
		return
	}

	r := routes.SetupRouter()
	// running
	r.Listen(":3000")
}
