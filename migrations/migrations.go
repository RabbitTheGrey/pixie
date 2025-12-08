package migrations

import (
	"pixie/lib/db/migrations"
)

func RegisterMigrations(mr migrations.IMigrationsResolver) {
	mr.AppendMigration("Version20251208215055", Version20251208215055_Up, Version20251208215055_Down)
	mr.AppendMigration("Version20251208215126", Version20251208215126_Up, Version20251208215126_Down)
}
