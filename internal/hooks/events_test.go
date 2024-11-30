package hooks

import (
	"testing"
)

func TestCastToString(t *testing.T) {
	// Example data
	event := constructEvent(orderItemEvent{
		OrderItemId: "12345",
		Status:      orderItemStatusInArbeit,
	})

	content := event.stringifyContent()

	// Expected JSON string
	expected := `{"order_item_id":"12345","status":"InArbeit"}`

	// Assert that the content matches the expected string
	if content != expected {
		t.Errorf("Content does not match expected JSON.\nGot: %s\nWant: %s", content, expected)
	}
}
