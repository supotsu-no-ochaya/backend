package hooks

import (
	"fmt"
	"github.com/pocketbase/dbx"
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

func requiresOrderItemStatusUpdateCheck(status orderStatus) bool {
	switch status {
	case orderStatusAufgegeben,
		orderStatusInArbeit,
		orderStatusAbholbereit,
		orderStatusGeliefert,
		orderStatusBezahlt:
		return true
	default:
		return false
	}
}

func mapOrderStatusToOrderItemStatus(orderStatus orderStatus) orderItemStatus {
	switch orderStatus {
	case orderStatusAufgegeben:
		return orderItemStatusAufgegeben
	case orderStatusInArbeit:
		return orderItemStatusInArbeit
	case orderStatusAbholbereit:
		return orderItemStatusAbholbereit
	case orderStatusGeliefert:
		return orderItemStatusGeliefert
	case orderStatusBezahlt:
		return orderItemStatusBezahlt
	default:
		panic("invalid order status") // TODO remove panic for error propagation
	}
}

func RegisterOrderHooks(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess(orderTableName).BindFunc(orderAfterCreateSuccess)
	app.OnRecordAfterUpdateSuccess(orderTableName).BindFunc(orderAfterUpdateSuccess)
}

func orderAfterCreateSuccess(orderRecordEvent *core.RecordEvent) error {
	orderEvent := orderEvent{
		OrderId: orderRecordEvent.Record.Get("id").(string),
		Status:  orderStatus(orderRecordEvent.Record.Get("status").(string)),
	}
	return constructEvent(orderEvent).save(orderRecordEvent.App)
}

func orderAfterUpdateSuccess(orderRecordEvent *core.RecordEvent) error {
	oldStatus := orderStatus(orderRecordEvent.Record.Original().GetString("status"))
	newStatus := orderStatus(orderRecordEvent.Record.GetString("status"))

	// If Status hasn't changed, no action is needed.
	if oldStatus == newStatus {
		return nil
	}

	app := orderRecordEvent.App
	orderEvent := orderEvent{
		OrderId: orderRecordEvent.Record.Get("id").(string),
		Status:  newStatus,
	}
	orderID := orderRecordEvent.Record.GetString("id")
	status := orderStatus(orderRecordEvent.Record.GetString("status"))
	app.Logger().Info(
		fmt.Sprintf("Order with id: %s changed to status %s ", orderID, status),
	)
	// Create an event record for the order item Status change.
	if err := constructEvent(orderEvent).save(orderRecordEvent.App); err != nil {
		return err
	}

	// find the "order" the updated "order item" belongs to
	// if all "order items" attached to that order are now in the same orderItemStatus set the order status to the equivilant status
	// e.g. if all order items are in status "InArbeit" set the order status to the "InArbeit" status as well.
	if requiresOrderItemStatusUpdateCheck(status) {
		app.Logger().Info(
			fmt.Sprintf("Updating order items stati because order %s changed into status %s", orderID, status),
		)
		orderItems, err := orderRecordEvent.App.FindRecordsByFilter(
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
		newOrderItemStatus := mapOrderStatusToOrderItemStatus(status)
		for _, orderItem := range orderItems {
			orderItem.Set("status", string(mapOrderItemStatusToOrderStatus(newOrderItemStatus)))
			err := orderRecordEvent.App.Save(orderItem)
			if err != nil {
				app.Logger().Error(
					fmt.Sprintf("failed to update order item with id: %s to status: %s", orderItem.Id, newOrderItemStatus),
				)
			}
		}

	}

	return nil
}
