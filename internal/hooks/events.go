package hooks

import (
	"encoding/json"
	"github.com/pocketbase/pocketbase/core"
)

const (
	eventTableName string = "event"
)

type event[T eventMapping] struct {
	eventType eventType
	content   T
}

type eventType string

const (
	orderEventType     = eventType(orderTableName)
	orderItemEventType = eventType(orderItemTableName)
	productEventType   = eventType(productTableName) // Added for product events
)

// Define the mapping between eventType and eventContent
type eventMapping interface {
	getEventType() eventType
}

// Order event
type orderEvent struct {
	OrderId string      `json:"order_id"`
	Status  orderStatus `json:"status"`
}

// Order item event
type orderItemEvent struct {
	OrderItemId string          `json:"order_item_id"`
	Status      orderItemStatus `json:"status"`
}

// Product event
type productEvent struct {
	ProductId   string `json:"product_id"`
	IsAvailable bool   `json:"is_available"`
}

// Associate `orderEvent` with `orderEventType`
func (orderEvent) getEventType() eventType {
	return orderEventType
}

// Associate `orderItemEvent` with `orderItemEventType`
func (orderItemEvent) getEventType() eventType {
	return orderItemEventType
}

// Associate `productEvent` with `productEventType`
func (productEvent) getEventType() eventType {
	return productEventType
}

func constructEvent[T eventMapping](content T) event[T] {
	return event[T]{
		eventType: content.getEventType(),
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
	record.Set("type", string(e.eventType))
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
