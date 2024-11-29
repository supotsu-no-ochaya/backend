package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"github.com/supotsu-no-ochaya/backend/internal/routes"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Initialize routes
		routes.RegisterAPIRoutes(e, app)
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
