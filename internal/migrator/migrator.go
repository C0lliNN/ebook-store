package migrator

import (
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DatabaseURL string
type Source string

type Config struct {
	DatabaseURL DatabaseURL
	Source      Source
}

type Migrator struct {
	Config
}

func New(c Config) *Migrator {
	return &Migrator{Config: c}
}

// Sync Applies new database migrations
func (m *Migrator) Sync() {
	mi, err := migrate.New(string(m.Source), string(m.DatabaseURL))

	if err != nil {
		log.Default().Fatalf("An error happened when trying to sync migrations: %v", err)
	}

	if err = mi.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Default().Debug("The current migrations are up to date")
			return
		}

		log.Default().Fatalf("An error happened when trying to sync migrations: %v", err)
	}
}