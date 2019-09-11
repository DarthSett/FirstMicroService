package router

import (
	"database/sql"
	"github.com/FirstMicroservice/CsvService/CsvUploader"
	"github.com/gin-gonic/gin"
)

type Router struct {
	db *sql.DB
}

func NewRouter(db *sql.DB) *Router {
	return &Router{db: db}
}

func (r *Router) Router() {
	s := gin.Default()
	s.GET("/status", r.status)
	s.GET("/archive",r.archive)
	s.POST("/upload", CsvUploader.Csvupload)
	s.Run(":4000")
}

func (r *Router) status(c *gin.Context) {
	row, err := r.db.Query("select link,status from Link")
	if err != nil {
		CsvUploader.FailOnError(err, "Failed to fetch the status from db: ")
	}
	m := make(map[string]string)
	defer row.Close()
	var (
		link   string
		status string
	)
	for row.Next() {
		err := row.Scan(&link, &status)
		CsvUploader.FailOnError(err, "Failed to scan the values")
		m[link] = status
	}
	c.JSON(200, m)
}

func (r *Router) archive(c *gin.Context) {
	row, err := r.db.Query("select link,archived from Link")
	if err != nil {
		CsvUploader.FailOnError(err, "Failed to fetch the archive info from db: ")
	}
	m := make(map[string]int)
	defer row.Close()
	var (
		link   string
		archived int
	)
	for row.Next() {
		err := row.Scan(&link, &archived)
		CsvUploader.FailOnError(err, "Failed to scan the values")
		m[link] = archived
	}
	c.JSON(200, m)
}