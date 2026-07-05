package migration

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dsn string) error {
	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Tidak ada migration baru, database sudah up to date")
			return nil
		}
		return err
	}

	log.Println("Migration berhasil dijalankan")
	return nil
}
