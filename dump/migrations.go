package dump

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"log"
	"os"
)

func MigrateDatabase(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/dump", dir),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}
	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
