package pgdb

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDbParams struct {
	Host     string
	User     string
	Password string
	Db       string
	Port     string
	Ssl      string
	Timezone string
}

type ReactRolesDatabase struct {
	DB *gorm.DB
}

func getConnectionString(params PostgresDbParams) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		params.Host, params.Port, params.User, params.Password, params.Db, params.Ssl, params.Timezone)
}

func runAutoMigrations(db *gorm.DB) {
	// Migrate tables
	db.AutoMigrate(&Role{}, &Selector{}, &ServerConfiguration{})
}

func GetDatabase(params PostgresDbParams) *ReactRolesDatabase {
	database, err := gorm.Open(postgres.Open(getConnectionString(params)), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("[PGDB] Connected to Postgres successfully")
	}

	runAutoMigrations(database)

	return &ReactRolesDatabase{
		DB: database,
	}
}
