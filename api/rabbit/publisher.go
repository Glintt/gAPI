package rabbit

import (
	"gAPIManagement/api/utils"
	"time"
	"github.com/streadway/amqp"
)

var RabbitConnectionGlobal *amqp.Connection
var RabbitChannelGlobal *amqp.Channel

func StartPublisherWorker() {
	var RabbitConnection *amqp.Connection
	var RabbitChannel *amqp.Channel
	var RabbitCloseError chan *amqp.Error
 	RabbitCloseError = make(chan *amqp.Error)
	RabbitConnection, RabbitChannel, RabbitCloseError = RabbitMQCreateConnection(RabbitConnection, RabbitChannel, RabbitCloseError)
	DeclareQueue(RabbitChannel)

	go RabbitMQConnectionReconnectHandle(RabbitConnection, RabbitChannel, RabbitCloseError)
	
	for {
		log := <-WorkerPool

		WorkerPublishToRabbitMQ(RabbitChannel,log)
	}
}

func StartPublisher() {
	RabbitConnectionGlobal = ConnectToRabbit()
    RabbitChannelGlobal = CreateChannel(RabbitConnectionGlobal)
    DeclareQueue(RabbitChannelGlobal)
}

func RabbitMQCreateConnection(RabbitConnection *amqp.Connection,RabbitChannel *amqp.Channel, RabbitCloseError chan *amqp.Error) (*amqp.Connection,*amqp.Channel, chan *amqp.Error){
	utils.LogMessage("Connect to RabbitMQ")
	RabbitConnection = ConnectToRabbit()
	utils.LogMessage("Connecteded to RabbitMQ")
	RabbitChannel = CreateChannel(RabbitConnection)
	utils.LogMessage("Channel created on connection to RabbitMQ")
	RabbitCloseError = make(chan *amqp.Error)
	RabbitConnection.NotifyClose(RabbitCloseError)

	return RabbitConnection, RabbitChannel, RabbitCloseError
}

func RabbitMQConnectionReconnectHandle(RabbitConnection *amqp.Connection,RabbitChannel *amqp.Channel, RabbitCloseError chan *amqp.Error) {
  var rabbitErr *amqp.Error

  for {
	rabbitErr = <-RabbitCloseError
    if rabbitErr != nil {
		RabbitMQCreateConnection(RabbitConnection, RabbitChannel, RabbitCloseError)
    }
  }
}

func ConnectToRabbit() *amqp.Connection {
	connectionString := "amqp://" + User() +":" + Pwd() + "@" + Host() + ":" + Port() 	+ "/"
	// InitPublishers(2)
	utils.LogMessage("Connecting to: " + connectionString)
	
	for {
		connection, err := amqp.Dial(connectionString)
		FailOnError(err, "Failed to connect to RabbitMQ")
		if err == nil {
			utils.LogMessage("Connected to " + connectionString)
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

func DeclareQueue(RabbitChannel *amqp.Channel) error{
	_, err := RabbitChannel.QueueDeclare(
		Queue(), // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          	// no-wait
		nil,            // arguments
	)
	if FailOnError(err, "Failed to declare queue") != nil {
		return err
	}
	return nil
}


func WorkerPublishToRabbitMQ(RabbitChannel *amqp.Channel, log []byte) error {
	err := RabbitChannel.Publish(
		"",            // exchange
		Queue(), // routing key
		false,         // mandatory
		false,         	// immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        log,
		})

	if FailOnError(err, "Failed to publish message") != nil {
		return err
	}
	return nil
}