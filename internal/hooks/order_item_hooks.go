package hooks

import (
	"fmt"
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
	orderItemStatusBezahlt     orderItemStatus = "Bezahlt"     //nolint:unused
)

func requiresOrderStatusUpdateCheck(status orderItemStatus) bool {
	switch status {
	case orderItemStatusAufgegeben,
		orderItemStatusInArbeit,
		orderItemStatusAbholbereit,
		orderItemStatusGeliefert,
		orderItemStatusBezahlt:
		return true
	default:
		return false
	}
}

func mapOrderItemStatusToOrderStatus(orderItemStatus orderItemStatus) orderStatus {
	switch orderItemStatus {
	case orderItemStatusAufgegeben:
		return orderStatusAufgegeben
	case orderItemStatusInArbeit:
		return orderStatusInArbeit
	case orderItemStatusAbholbereit:
		return orderStatusAbholbereit
	case orderItemStatusGeliefert:
		return orderStatusGeliefert
	case orderItemStatusBezahlt:
		return orderStatusBezahlt
	default:
		panic("invalid order item status") // TODO remove panic for error propagation
	}
}

func RegisterOrderItemHooks(app *pocketbase.PocketBase) {
	app.OnRecordAfterCreateSuccess(orderItemTableName).BindFunc(orderItemAfterCreateSuccess)
	app.OnRecordAfterUpdateSuccess(orderItemTableName).BindFunc(orderItemAfterUpdateSuccess)
}

func orderItemAfterCreateSuccess(orderItemRecordEvent *core.RecordEvent) error {
	orderItemEvent := orderItemEvent{
		OrderItemId: orderItemRecordEvent.Record.Get("id").(string),
		Status:      orderItemStatus(orderItemRecordEvent.Record.Get("status").(string)),
	}

	return constructEvent(orderItemEvent).save(orderItemRecordEvent.App)
}

func orderItemAfterUpdateSuccess(orderItemRecordEvent *core.RecordEvent) error {
	oldStatus := orderItemStatus(orderItemRecordEvent.Record.Original().GetString("status"))
	newStatus := orderItemStatus(orderItemRecordEvent.Record.GetString("status"))

	// If Status hasn't changed, no action is needed.
	if oldStatus == newStatus {
		return nil
	}
	return handleOrderItemStatusUpdate(orderItemRecordEvent)
}

func handleOrderItemStatusUpdate(orderItemRecordEvent *core.RecordEvent) error {
	app := orderItemRecordEvent.App
	orderItemEvent := orderItemEvent{
		OrderItemId: orderItemRecordEvent.Record.Get("id").(string),
		Status:      orderItemStatus(orderItemRecordEvent.Record.Get("status").(string)),
	}
	// Create an event record for the order item Status change.
	if err := constructEvent(orderItemEvent).save(orderItemRecordEvent.App); err != nil {
		return err
	}

	status := orderItemStatus(orderItemRecordEvent.Record.GetString("status"))
	// find the "order" the updated "order item" belongs to
	// if all "order items" attached to that order are now in the same orderItemStatus set the order status to the equivilant status
	// e.g. if all order items are in status "InArbeit" set the order status to the "InArbeit" status as well.
	if requiresOrderStatusUpdateCheck(status) {
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
		if allOrderItemsHaveStatus(orderItems, orderItemRecordEvent.Record.GetString("status")) {
			app.Logger().Info(
				fmt.Sprintf("All order items of order (id: %s) are in status: %s ... Updating order status.", orderID, status),
			)
			order, err := orderItemRecordEvent.App.FindRecordById("order", orderID)
			if err != nil {
				app.Logger().Error(
					fmt.Sprintf("Failed to find order with id: %s", orderID),
				)
				return err
			}
			order.Set("status", string(mapOrderItemStatusToOrderStatus(status)))
			orderUpdateErr := orderItemRecordEvent.App.Save(order)
			if orderUpdateErr != nil {
				app.Logger().Error(
					fmt.Sprintf("Failed to save order with id: %s", orderID),
				)
				return orderUpdateErr
			}
			app.Logger().Info(
				fmt.Sprintf("Successfully updated order with id: %s to status: %s", orderID, status),
			)

			return nil
		}
	}
	return nil
}

func allOrderItemsHaveStatus(orderItems []*core.Record, status string) bool {
	for _, item := range orderItems {
		if item.GetString("status") != status {
			return false
		}
	}
	return true
}
