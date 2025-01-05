package database

import (
	"os"

	"github.com/rs/zerolog/log"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Database struct {
	*gorm.DB
}

var (
	dbUrl      = os.Getenv("DB_URL")
	dbInstance *Database
)

func NewInstance() *Database {
	gorm, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        dbUrl,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// Logger: gorm_zerolog.New(),
	})

	if err != nil {
		log.Fatal().Err(err)
	}

	return &Database{gorm}
}

func Instance() *Database {
	if dbInstance != nil {
		return dbInstance
	}

	dbInstance = NewInstance()
	return dbInstance
}

func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	log.Info().Msg("Closed database connection")
	return sqlDB.Close()
}
