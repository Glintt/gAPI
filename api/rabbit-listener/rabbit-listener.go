package main

import (
	"encoding/json"
	"os"
	"strconv"

	apianalytics "github.com/Glintt/gAPI/api/api-analytics"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/database"
	"github.com/Glintt/gAPI/api/logs"
	logsModels "github.com/Glintt/gAPI/api/logs/models"
	"github.com/Glintt/gAPI/api/rabbit"
	"github.com/Glintt/gAPI/api/utils"

	"github.com/streadway/amqp"
)

var ELASTIC_URL string
var ELASTICPORT string

func teste() {
	config.LoadURLConstants()
	config.LoadGApiConfig()

	if config.GApiConfiguration.Logs.Type == apianalytics.LogOracleType {
		database.InitDatabaseConnection()
	}

	workers := 1
	if os.Getenv("RABBIT_LISTENER_WORKERS") != "" {
		workers, _ = strconv.Atoi(os.Getenv("RABBIT_LISTENER_WORKERS"))
		Start(workers)
	} else {
		StartListeningToRabbit(1)
	}
}

func Start(workers int) {
	/* if (os.Getenv("ELASTICSEARCH_HOST") != "") {
	    ELASTIC_URL = os.Getenv("ELASTICSEARCH_HOST")
	} */

	forever := make(chan bool)
	for i := 0; i < workers; i++ {
		go StartListeningToRabbit(i)
	}
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		utils.LogMessage(msg+err.Error(), utils.ErrorLogType)
	}
}

func PreventCrash() {
	if r := recover(); r != nil {
		utils.LogMessage("Rabbit Listener Crashed", utils.ErrorLogType)
		StartListeningToRabbit(1)
	}
}

func StartListeningToRabbit(workerId int) {
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
		true,           // durable
		false,          // delete when usused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
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

	go ReceiveAndPublish(workerId, msgs)

	utils.LogMessage(" [*] Waiting for messages. To exit press CTRL+C", utils.InfoLogType)
	<-forever
}

func ReceiveAndPublish(workerId int, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		var reqLogging logsModels.RequestLogging
		err := json.Unmarshal(d.Body, &reqLogging)
		if err == nil {
			utils.LogMessage("Publish to elasticsearch from #"+strconv.Itoa(workerId)+" - "+string(d.Body), utils.InfoLogType)
			// logs.LoggingType[config.GApiConfiguration.Logs.Type].(func(*logsModels.RequestLogging))(&reqLogging)
			logs.LoggingType[config.GApiConfiguration.Logs.Type].(func(*logsModels.RequestLogging))(&reqLogging)

		} else {
			utils.LogMessage("Error logging message: "+string(d.Body), utils.ErrorLogType)
		}
	}
}
