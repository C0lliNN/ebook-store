package db

import (
	"database/sql"
	"github.com/c0llinn/ebook-store/config/log"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection() *gorm.DB {
	db, err := sql.Open("postgres", viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Logger.Fatalw("Postgres connection has failed", "error", err.Error())
		return nil
	}

	dialector := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Logger.Fatalw("Postgres connection has failed", "error", err.Error())
		return nil
	}

	if err = db.Ping(); err != nil {
		log.Logger.Fatalw("Ping has failed", "error", err.Error())
		return nil
	}

	return conn
}
