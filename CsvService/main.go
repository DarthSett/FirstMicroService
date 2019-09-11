package main

import (
	"database/sql"
	"fmt"
	"github.com/FirstMicroservice/CsvService/CsvUploader"
	"github.com/FirstMicroservice/CsvService/router"
)

func main() {

	con := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", "root", "password", "database", "3306", "mslinks")
	db, err := sql.Open("mysql", con)
	CsvUploader.FailOnError(err,"Can't connect to db")
	defer db.Close()
	r := router.NewRouter(db)
	r.Router()
}
