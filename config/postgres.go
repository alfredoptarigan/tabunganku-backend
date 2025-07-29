package config

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

var (
	databasePostgresOnce     sync.Once
	databasePostgresInstance *gorm.DB
)

type DatabasePostgres Database

func InitDatabasePostgres() *gorm.DB {
	// This function initializes the database connection
	databasePostgresOnce.Do(func() {
		databasePostgresInstance = &gorm.DB{}

		// Initialize the database connection here
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			DbHost,
			DbUser,
			DbPassword,
			DbName,
			DbPort,
			DbSSLMode,
			DbTimezone,
		)

		// println all of the config
		fmt.Println("dsn:", dsn)

		// Open the database connection
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic(err)
		}

		// Set the connection pool settings
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get database connection pool")
		}

		// Set the maximum number of open connections
		// Set the maximum number of open connections
		dbMaxConnection, _ := strconv.ParseInt(DbMaxConnections, 10, 64)
		dbMaxIdleConnection, _ := strconv.ParseInt(DbIdleConnections, 10, 64)
		sqlDB.SetMaxOpenConns(int(dbMaxConnection))
		// Set the maximum number of idle connections
		sqlDB.SetMaxIdleConns(int(dbMaxIdleConnection))
		// Set the maximum lifetime of a connection
		sqlDB.SetConnMaxLifetime(1 * time.Hour)

		databasePostgresInstance = db
	})

	return databasePostgresInstance
}
