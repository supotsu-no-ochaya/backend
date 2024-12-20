package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		paymentOptions, err := app.FindCollectionByNameOrId("payment_option")
		if err != nil {
			return err
		}

		options := []string{"Bar", "Karte"}
		for _, option := range options {
			record := core.NewRecord(paymentOptions)
			record.Set("name", option)

			if err := app.Save(record); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		options := []string{"Bar", "Karte"}
		for _, option := range options {
			record, err := app.FindFirstRecordByData("payment_option", "name", option)
			if err != nil {
				return err
			}
			if record != nil {
				if err := app.Delete(record); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
