package hooks

import (
	"github.com/pocketbase/pocketbase"
)

// RegisterHooks sets up the event listeners for record updates and creations.
func RegisterHooks(app *pocketbase.PocketBase) {
	registerOrderHooks(app)
	registerOrderItemHooks(app)
}
