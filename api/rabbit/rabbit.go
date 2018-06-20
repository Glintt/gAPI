package rabbit

import (
	"gAPIManagement/api/utils"	
	"os"
)

var WorkerPool chan []byte

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

func FailOnError(err error, msg string) error {
	if err != nil {
		utils.LogMessage(msg + ": " + err.Error())
	}
	return err
}

func InitPublishers(workers int) {
	WorkerPool = make(chan []byte)
	for i:=0; i < workers; i++ {
		go StartPublisherWorker()
	}
}

func PublishToRabbitMQ(log []byte){
	WorkerPool <- log
}