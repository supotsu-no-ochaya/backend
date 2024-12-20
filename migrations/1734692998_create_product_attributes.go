package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		productAttribute, err := app.FindCollectionByNameOrId("product_attribute")
		if err != nil {
			return err
		}

		attributes := []string{"veggie", "vegan"}
		for _, attribute := range attributes {
			record := core.NewRecord(productAttribute)
			record.Set("name", attribute)

			if err := app.Save(record); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		attributes := []string{"veggie", "vegan"}
		for _, attribute := range attributes {
			record, err := app.FindFirstRecordByData("product_attribute", "name", attribute)
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
