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

func orderAfterCreateSuccess(orderRecordEvent *core.RecordEvent) error {
	orderEvent := orderEvent{
		OrderId: orderRecordEvent.Record.Get("id").(string),
		Status:  orderRecordEvent.Record.Get("status").(orderStatus),
	}
	return constructEvent(orderEvent).save(orderRecordEvent.App)
}

func orderAfterUpdateSuccess(orderRecordEvent *core.RecordEvent) error {
	oldStatus := orderStatus(orderRecordEvent.Record.Original().GetString("status"))
	newStatus := orderStatus(orderRecordEvent.Record.GetString("status"))

	// If Status hasn't changed, no action is needed.
	if oldStatus != newStatus {
		orderEvent := orderEvent{
			OrderId: orderRecordEvent.Record.Get("id").(string),
			Status:  orderRecordEvent.Record.Get("status").(orderStatus),
		}
		// Create an event record for the order item Status change.
		if err := constructEvent(orderEvent).save(orderRecordEvent.App); err != nil {
			return err
		}
	}

	return nil
}
