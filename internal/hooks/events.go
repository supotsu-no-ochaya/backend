package hooks

import (
	"encoding/json"
	"github.com/pocketbase/pocketbase/core"
)

const (
	eventTableName string = "event"
)

type event[T eventContent] struct {
	eventType eventType
	content   T
}

type eventType string

const (
	orderEventType     = eventType(orderTableName)
	orderItemEventType = eventType(orderItemTableName)
)

type eventContent interface {
	orderItemEvent | orderEvent
}

type orderEvent struct {
	OrderId string      `json:"order_id"` // JSON key will be "order_id"
	Status  orderStatus `json:"Status"`   // JSON key will be "Status"
}

type orderItemEvent struct {
	OrderItemId string          `json:"order_item_id"` // JSON key will be "order_item_id"
	Status      orderItemStatus `json:"Status"`        // JSON key will be "Status"
}

func constructEvent[T eventContent](content T) event[T] {
	var eventType eventType

	// Determine the eventType based on the type of content
	switch any(content).(type) {
	case orderItemEvent:
		eventType = orderItemEventType
	case orderEvent:
		eventType = orderEventType
	default:
		panic("Unsupported event content type")
	}

	return event[T]{
		eventType: eventType,
		content:   content,
	}
}

func (e event[T]) save(app core.App) error {
	collection, err := app.FindCollectionByNameOrId(eventTableName)
	if err != nil {
		return err
	}

	contentString := e.stringifyContent()
	record := core.NewRecord(collection)
	record.Set("type", e.eventType)
	record.Set("content", contentString)

	return app.Save(record)
}

func (e *event[T]) stringifyContent() string {
	data, err := json.Marshal(e.content)
	if err != nil {
		panic("Cannot marshal event")
	}
	return string(data)
}
