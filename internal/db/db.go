package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
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

func Get() *Database {
	if dbInstance != nil {
		return dbInstance
	}

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
		log.Fatal().Err(err).Msg("")
	}

	dbInstance = &Database{gorm}

	return dbInstance
}

func (db *Database) Health() map[string]string {
	stats := make(map[string]string)

	sqlDB, err := db.DB.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatal().Err(err).Msg("db down")
		return stats
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatal().Err(err).Msg("db down")
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	// ...existing code for stats evaluation...

	return stats
}

func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	log.Printf("Disconnected from database: %s", dbUrl)
	return sqlDB.Close()
}
