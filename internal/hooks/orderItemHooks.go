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

func orderItemAfterCreateSuccess(orderItemRecordEvent *core.RecordEvent) error {
	orderItemEvent := orderItemEvent{
		OrderItemId: orderItemRecordEvent.Record.Get("id").(string),
		Status:      orderItemRecordEvent.Record.Get("Status").(orderItemStatus),
	}

	return constructEvent(orderItemEvent).save(orderItemRecordEvent.App)
}

func orderItemAfterUpdateSuccess(orderItemRecordEvent *core.RecordEvent) error {
	oldStatus := orderItemStatus(orderItemRecordEvent.Record.Original().GetString("Status"))
	newStatus := orderItemStatus(orderItemRecordEvent.Record.GetString("Status"))

	// If Status hasn't changed, no action is needed.
	if oldStatus != newStatus {
		return handleOrderItemStatusUpdate(orderItemRecordEvent)
	}

	return nil
}

func handleOrderItemStatusUpdate(orderItemRecordEvent *core.RecordEvent) error {
	orderItemEvent := orderItemEvent{
		OrderItemId: orderItemRecordEvent.Record.Get("id").(string),
		Status:      orderItemRecordEvent.Record.Get("Status").(orderItemStatus),
	}
	// Create an event record for the order item Status change.
	if err := constructEvent(orderItemEvent).save(orderItemRecordEvent.App); err != nil {
		return err
	}

	// find the "order" the updated "order item" belongs to
	// if all "order items" attached to that order are now in Status "abholbereit, set the "order" Status to "abholbereit"
	if orderItemStatus(orderItemRecordEvent.Record.GetString("Status")) == orderItemStatusInArbeit {
		orderID := orderItemRecordEvent.Record.GetString("order")
		orderItems, err := orderItemRecordEvent.App.FindRecordsByFilter(
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
		if allOrderItemsHaveStatus(orderItems, orderItemRecordEvent.Record.GetString("Status")) {
			order, err := orderItemRecordEvent.App.FindRecordById("order", orderID)
			if err != nil {
				return err
			}
			order.Set("Status", orderStatusInArbeit)
			return orderItemRecordEvent.App.Save(order)
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
