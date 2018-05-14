package rabbit

import (
	"gAPIManagement/api/utils"
	"time"
	"github.com/streadway/amqp"
	
	"os"
)

func User() string{
	return os.Getenv("RABBITMQ_USER")
}

func Pwd() string{
	return os.Getenv("RABBITMQ_PASSWORD")
}

func Host() string{
	return os.Getenv("RABBITMQ_HOST")
}

func Port() string{
	return os.Getenv("RABBITMQ_PORT")
}

func Queue() string{
	return os.Getenv("RABBITMQ_QUEUE")
}


func ConnectToRabbit() *amqp.Connection {
	utils.LogMessage("amqp://" + User() +":" + Pwd() + "@" + Host() + ":" + Port() 	+ "/")
	
	for {
		connection, err := amqp.Dial("amqp://" + User() +":" + Pwd() + "@" + Host() + ":" + Port() 	+ "/")
		FailOnError(err, "Failed to connect to RabbitMQ")
		if err == nil {
			return connection
		}
		
		time.Sleep(500 * time.Millisecond)
	}
}

func CreateChannel(connection *amqp.Connection) *amqp.Channel {
	rc, err := connection.Channel()
	FailOnError(err, "Failed to open a channel")
	return rc
}

func FailOnError(err error, msg string) {
	if err != nil {
	  utils.LogMessage( msg + ":" + err.Error())
	}
}