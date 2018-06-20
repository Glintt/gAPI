package main

import (
	"gAPIManagement/api/logs"
	"gAPIManagement/api/utils"
	"gAPIManagement/api/rabbit"
	"os"
	"encoding/json"
)

var ELASTIC_URL string
var ELASTICPORT string

func main(){
	StartListeningToRabbit()
}


func Start(workers int) {
    if (os.Getenv("ELASTICSEARCH_HOST") != "") {
        ELASTIC_URL = os.Getenv("ELASTICSEARCH_HOST")
    }

	for i:=0; i < workers; i++ {
		go StartListeningToRabbit()
	}
}


func failOnError(err error, msg string) {
  if err != nil {
   utils.LogMessage(msg + err.Error())
  }
}


func PreventCrash(){
	if r := recover(); r != nil {
		utils.LogMessage("Rabbit Listener Crashed")
		StartListeningToRabbit()
	}
}

func StartListeningToRabbit() {
	defer PreventCrash()

	ELASTIC_URL = os.Getenv("ELASTICSEARCH_HOST")
	ELASTICPORT = os.Getenv("ELASTICSEARCH_PORT")
	
	conn := rabbit.ConnectToRabbit()
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel 2")
	defer ch.Close()

	q, err := ch.QueueDeclare(
	rabbit.Queue(), // name
	true,   // durable
	false,   // delete when usused
	false,   // exclusive
	false,   // no-wait
	nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue 2")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	failOnError(err, "Failed to register a consumer 2")
	  
	forever := make(chan bool)
	  
	go func() {
		for d := range msgs {
			
			var reqLogging logs.RequestLogging
			err := json.Unmarshal(d.Body, &reqLogging)
			if err == nil{				
				utils.LogMessage("Publish to elasticsearch")
				logs.PublishElastic(&reqLogging)
			}else{
				utils.LogMessage("Error logging message: " + string(d.Body))
			}
		}
	}()
	  
	utils.LogMessage(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
