package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/supotsu-no-ochaya/backend/internal/hooks"
	"github.com/supotsu-no-ochaya/backend/internal/routes"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		routes.RegisterAPIRoutes(e, app)
		return e.Next()
	})

	hooks.RegisterOrderHooks(app)
	hooks.RegisterOrderItemHooks(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
