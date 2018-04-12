package main

import (
	"gAPIManagement/api/config"
	"gAPIManagement/api/utils"
	"fmt"
	"gAPIManagement/api/rabbit"
	"gAPIManagement/api/logs"
	"os"
	"encoding/json"
	)

func failOnError(err error, msg string) {
  if err != nil {
    fmt.Println("%s: %s", msg, err)
  }
}

var ELASTICURL string
var ELASTICPORT string

func StartListeningToRabbit() {
	utils.PreventCrash()

	ELASTICURL = os.Getenv("ELASTICSEARCH_HOST")
	ELASTICPORT = os.Getenv("ELASTICSEARCH_PORT")
	
	conn := rabbit.ConnectToRabbit()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
	rabbit.Queue(), // name
	true,   // durable
	false,   // delete when usused
	false,   // exclusive
	false,   // no-wait
	nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	failOnError(err, "Failed to register a consumer")
	  
	forever := make(chan bool)
	  
	go func() {
		for d := range msgs {
			var reqLogging logs.RequestLogging
			err := json.Unmarshal(d.Body, &reqLogging)
			if err == nil{
				logs.PublishElastic(&reqLogging)
			}else{
				fmt.Printf("Error logging message: %s", d.Body)		
			}
		}
	}()
	  
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main(){
	config.LoadURLConstants()
	
	StartListeningToRabbit()
}

