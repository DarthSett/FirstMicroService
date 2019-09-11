package main

import (
	"database/sql"
	"fmt"
	scraper "github.com/FirstMicroservice/ScraperService/Scraper"
	"github.com/FirstMicroservice/dump"
	"github.com/streadway/amqp"
	"log"
	"time"
)

//func failOnError(err error, msg string) {
//	if err != nil {
//		log.Fatalf("%s: %s", msg, err)
//	}
//}

func main() {
	con := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", "root", "password", "database", "3306", "mslinks")
	db, err := sql.Open("mysql", con)
	scraper.FailOnError(err, "Can't connect to db")
	println("@@@@ DB Connected")
	defer db.Close()

	println("calling rconnect")
	conn := rConnect()
	println("leaving rconnect")
	defer conn.Close()

	ch, err := conn.Channel()
	scraper.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	scraper.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	scraper.FailOnError(err, "Failed to register a consumer")
	println("calling rmigrate")
	_ = rMigrate(db)
	if err == nil {
		println("db migrated")
	}
	scraper.FailOnError(err, "Can't migrate db")

	forever := make(chan string)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			if string(d.Body) == "CSV Uploaded" {
				scraper.Scrape(db)
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func rConnect() *amqp.Connection {

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		println(err.Error())
		println("trying to rc")
		time.Sleep(5 * time.Second)
		return rConnect()
	}

	return conn
}

func rMigrate(db *sql.DB) int {

	err := dump.MigrateDatabase(db)
	if err != nil {
		println(err.Error())
		println("trying to reMigrate")
		time.Sleep(5 * time.Second)
		return rMigrate(db)
	}

	return 0
}
