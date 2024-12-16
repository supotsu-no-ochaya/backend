package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"os"
)

func init() {
	m.Register(func(app core.App) error {
		superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		record := core.NewRecord(superusers)

		// Load email and password from environment variables with default values
		email := os.Getenv("SUPERUSER_EMAIL")
		if email == "" {
			email = "admin@admin.admin"
		}

		password := os.Getenv("SUPERUSER_PASSWORD")
		if password == "" {
			password = "1234567890"
		}

		record.Set("email", email)
		record.Set("password", password)

		return app.Save(record)
	}, func(app core.App) error { // optional revert operation
		email := os.Getenv("SUPERUSER_EMAIL")
		if email == "" {
			email = "admin@admin.admin"
		}

		record, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, email)
		if record == nil {
			return nil
		}

		return app.Delete(record)
	})
}
