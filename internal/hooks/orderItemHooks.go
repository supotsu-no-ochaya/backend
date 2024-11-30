package hooks

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	orderItemTableName string = "order_item"
)

type orderItemStatus string

const (
	orderItemStatusAufgegeben  orderItemStatus = "Aufgegeben"  //nolint:unused
	orderItemStatusInArbeit    orderItemStatus = "InArbeit"    //nolint:unused
	orderItemStatusAbholbereit orderItemStatus = "Abholbereit" //nolint:unused
	orderItemStatusGeliefert   orderItemStatus = "Geliefert"   //nolint:unused
)

func RegisterOrderItemHooks(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess(orderItemTableName).BindFunc(orderItemAfterCreateSuccess)
	app.OnRecordAfterUpdateSuccess(orderItemTableName).BindFunc(orderItemAfterUpdateSuccess)
}

func orderItemAfterCreateSuccess(e *core.RecordEvent) error {
	orderItemEvent := orderItemEvent{
		OrderItemId: e.Record.Get("id").(string),
		Status:      e.Record.Get("Status").(orderItemStatus),
	}

	return constructEvent(orderItemEvent).save(e.App)
}

func orderItemAfterUpdateSuccess(e *core.RecordEvent) error {
	oldStatus := orderItemStatus(e.Record.Original().GetString("Status"))
	newStatus := orderItemStatus(e.Record.GetString("Status"))

	// If Status hasn't changed, no action is needed.
	if oldStatus != newStatus {
		return handleOrderItemStatusUpdate(e)
	}

	return nil
}

func handleOrderItemStatusUpdate(e *core.RecordEvent) error {
	orderItemEvent := orderItemEvent{
		OrderItemId: e.Record.Get("id").(string),
		Status:      e.Record.Get("Status").(orderItemStatus),
	}
	// Create an event record for the order item Status change.
	if err := constructEvent(orderItemEvent).save(e.App); err != nil {
		return err
	}

	// find the "order" the updated "order item" belongs to
	// if all "order items" attached to that order are now in Status "abholbereit, set the "order" Status to "abholbereit"
	if orderItemStatus(e.Record.GetString("Status")) == orderItemStatusInArbeit {
		orderID := e.Record.GetString("order")
		orderItems, err := e.App.FindRecordsByFilter(
			orderItemTableName,
			"order = {:OrderId}",
			"",
			0,
			0,
			dbx.Params{"OrderId": orderID},
		)
		if err != nil {
			return err
		}
		if allOrderItemsHaveStatus(orderItems, e.Record.GetString("Status")) {
			order, err := e.App.FindRecordById("order", orderID)
			if err != nil {
				return err
			}
			order.Set("Status", orderStatusInArbeit)
			return e.App.Save(order)
		}
	}
	return nil
}

func allOrderItemsHaveStatus(orderItems []*core.Record, status string) bool {
	for _, item := range orderItems {
		if item.GetString("Status") != status {
			return false
		}
	}
	return true
}
