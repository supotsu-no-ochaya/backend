package hooks

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterHooks sets up the event listeners for record updates and creations.
func RegisterHooks(app *pocketbase.PocketBase) {
	app.OnRecordAfterUpdateSuccess().BindFunc(handleAfterUpdateSuccess)
	app.OnRecordAfterCreateSuccess().BindFunc(handleAfterCreateSuccess)
}

// handleAfterUpdateSuccess handles post-update logic for specific collections.
func handleAfterUpdateSuccess(e *core.RecordEvent) error {
	collectionName := e.Record.Collection().Name

	switch collectionName {
	case "order", "order_item":
		return handleStatusUpdate(e)
	default:
		return nil
	}
}

// handleAfterCreateSuccess handles post-create logic for specific collections.
func handleAfterCreateSuccess(e *core.RecordEvent) error {
	collectionName := e.Record.Collection().Name

	if collectionName == "order" || collectionName == "order_item" {
		eventContent := buildEventContent(e.Record)
		return createEventRecord(e.App, collectionName, eventContent)
	}

	return nil
}

// handleStatusUpdate processes status changes for "order" and "order_item" collections.
func handleStatusUpdate(e *core.RecordEvent) error {
	oldStatus := e.Record.Original().GetString("status")
	newStatus := e.Record.GetString("status")

	// If status hasn't changed, no action is needed.
	if oldStatus == newStatus {
		return nil
	}

	collectionName := e.Record.Collection().Name
	eventContent := buildEventContent(e.Record)

	// Create an event record for the status change.
	if err := createEventRecord(e.App, collectionName, eventContent); err != nil {
		return err
	}

	// Additional logic for "order_item" collection.
	if collectionName == "order_item" {
		return handleOrderItemStatusChange(e)
	}

	return nil
}

// handleOrderItemStatusChange checks if all order items have the same status and updates the parent order accordingly.
func handleOrderItemStatusChange(e *core.RecordEvent) error {
	orderID := e.Record.GetString("order")
	orderItems, err := e.App.FindRecordsByFilter(
		"order_item",
		"order = {:orderId}",
		"",
		0,
		0,
		dbx.Params{"orderId": orderID},
	)
	if err != nil {
		return err
	}

	// Check if all order items have the same status.
	if allOrderItemsHaveStatus(orderItems, e.Record.GetString("status")) {
		order, err := e.App.FindRecordById("order", orderID)
		if err != nil {
			return err
		}

		// Update the order status only if the new status is "in arbeit".
		newStatus := e.Record.GetString("status")
		if newStatus != "in arbeit" {
			return nil
		}

		order.Set("status", newStatus)
		return e.App.Save(order)
	}

	return nil
}

// allOrderItemsHaveStatus checks if all order items have the specified status.
func allOrderItemsHaveStatus(orderItems []*core.Record, status string) bool {
	for _, item := range orderItems {
		if item.GetString("status") != status {
			return false
		}
	}
	return true
}

// buildEventContent constructs the content map for event records.
func buildEventContent(record *core.Record) map[string]interface{} {
	return map[string]interface{}{
		record.Collection().Name + "_id": record.Get("id"),
		"status":                         record.Get("status"),
	}
}

// createEventRecord creates a new event record in the "event" collection.
func createEventRecord(app core.App, eventType string, content map[string]interface{}) error {
	collection, err := app.FindCollectionByNameOrId("event")
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	record.Set("type", eventType)
	record.Set("content", content)

	return app.Save(record)
}
