package hooks

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	productTableName string = "product"
)

func RegisterProductHooks(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess().BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Collection().Name == productTableName {
			return productAfterCreateSuccess(e)
		}
		return e.Next()
	})

	app.OnRecordAfterUpdateSuccess().BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Collection().Name == productTableName {
			return productAfterUpdateSuccess(e)
		}
		return e.Next()
	})
}

func productAfterCreateSuccess(e *core.RecordEvent) error {
	// On creation, we can log the initial availability if needed
	productEvent := productEvent{
		ProductId:   e.Record.GetString("id"),
		IsAvailable: e.Record.GetBool("is_available"),
	}
	return constructEvent(productEvent).save(e.App)
}

func productAfterUpdateSuccess(e *core.RecordEvent) error {
	oldAvailable := e.Record.Original().GetBool("is_available")
	newAvailable := e.Record.GetBool("is_available")

	if oldAvailable != newAvailable {
		productEvent := productEvent{
			ProductId:   e.Record.GetString("id"),
			IsAvailable: newAvailable,
		}
		if err := constructEvent(productEvent).save(e.App); err != nil {
			return err
		}
	}

	return e.Next()
}
