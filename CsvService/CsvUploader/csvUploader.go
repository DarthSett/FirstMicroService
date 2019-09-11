package CsvUploader

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"io"
	"log"
	"time"
)

func Csvupload(c *gin.Context) {

	file, _, err := c.Request.FormFile("file")
	FailOnError( err,"There was a problem while reading file: ")
	defer file.Close()
	first := true
	c.Header("Content-type","text/plain")
	c.Writer.Write([]byte("file uploaded"))
	reader := csv.NewReader(file)
	var links []string
	for {
		record, err := reader.Read()
		if first {
			first = false
			continue
		}
		if err == io.EOF {
			break
		}
		links = append(links, record[0])
	}


	dbInsert(links)
	sendSignal("CSV Uploaded")

}

func dbInsert(links []string) {
	db := rConnect()
	for _,l := range links {
		q := fmt.Sprintf("INSERT into Link values ('%s','unscraped')",l)

		_,err := db.Query(q)
		FailOnError(err, "Failed to insert values into table")
		//if err != nil {
		//	e ,ok := err.(*mysql.MySQLError)
		//	if !ok {
		//		FailOnError(err,"didn't convert error: ")
		//	}
		//	if e.Number == 1062{
		//		continue
		//	}else{
		//		FailOnError(err, "Failed to insert values into table")
		//	}
		//}
	}
	db.Close()
}

func sendSignal(body string){
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()



	ch, err := conn.Channel()

	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode:amqp.Persistent,
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	FailOnError(err, "Failed to publish a message")
}

func FailOnError(err error, msg string) {
	if err != nil {
		println(msg + ": " + err.Error())
	}
}

func rConnect() *sql.DB{
	con := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", "root", "password", "database", "3306", "mslinks")
	db, err := sql.Open("mysql", con)
	if err!=nil{
		println(err.Error())
		println("trying to rc")
		time.Sleep(5 *time.Second)
		return rConnect()
	}

	return db
}