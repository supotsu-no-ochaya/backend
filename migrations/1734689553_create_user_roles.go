package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		userRoles, err := app.FindCollectionByNameOrId("user_role")
		if err != nil {
			return err
		}

		roles := []string{"Kuechenchef", "Kellner", "Kueche"}
		for _, role := range roles {
			record := core.NewRecord(userRoles)
			record.Set("role_name", role)

			if err := app.Save(record); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		roles := []string{"Kuechenchef", "Kellner", "Kueche"}
		for _, role := range roles {
			record, err := app.FindFirstRecordByData("user_role", "role_name", role)
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
