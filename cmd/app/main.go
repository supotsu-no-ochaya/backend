package main

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/supotsu-no-ochaya/backend/internal/hooks"
	"github.com/supotsu-no-ochaya/backend/internal/routes"
	_ "github.com/supotsu-no-ochaya/backend/migrations"
	"log"
)

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{})

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
