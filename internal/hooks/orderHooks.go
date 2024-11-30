package hooks

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	orderTableName string = "order"
)

type orderStatus string

const (
	orderStatusAufgegeben  orderStatus = "Aufgegeben"
	orderStatusInArbeit    orderStatus = "InArbeit"
	orderStatusAbholbereit orderStatus = "Abholbereit"
	orderStatusGeliefert   orderStatus = "Geliefert"
	orderStatusBezahlt     orderStatus = "Bezahlt"
)

func RegisterOrderHooks(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess(orderItemTableName).BindFunc(orderAfterCreateSuccess)
	app.OnRecordAfterUpdateSuccess(orderItemTableName).BindFunc(orderAfterUpdateSuccess)
}

func orderAfterCreateSuccess(e *core.RecordEvent) error {
	orderEvent := orderEvent{
		OrderId: e.Record.Get("id").(string),
		Status:  e.Record.Get("Status").(orderStatus),
	}
	return constructEvent(orderEvent).save(e.App)
}

func orderAfterUpdateSuccess(e *core.RecordEvent) error {
	oldStatus := orderStatus(e.Record.Original().GetString("Status"))
	newStatus := orderStatus(e.Record.GetString("Status"))

	// If Status hasn't changed, no action is needed.
	if oldStatus != newStatus {
		orderEvent := orderEvent{
			OrderId: e.Record.Get("id").(string),
			Status:  e.Record.Get("Status").(orderStatus),
		}
		// Create an event record for the order item Status change.
		if err := constructEvent(orderEvent).save(e.App); err != nil {
			return err
		}
	}

	return nil
}
