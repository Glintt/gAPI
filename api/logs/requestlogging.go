package logs

import (
	"encoding/json"
	"gAPIManagement/api/config"
	"gAPIManagement/api/logs/models"
	"gAPIManagement/api/logs/providers"
	"gAPIManagement/api/rabbit"
	"gAPIManagement/api/utils"

	// "github.com/streadway/amqp"

	"github.com/valyala/fasthttp"
)

var LoggingType = map[string]interface{}{
	"Rabbit":  PublishRabbit,
	"Elastic": providers.PublishElastic,
	"Oracle":  providers.PublishOracle}

func NewRequestLogging(c *fasthttp.RequestCtx, queryArgs []byte, headers []byte, currentDate string, elapsedTime int64, serviceName string, indexName string) models.RequestLogging {
	remoteAddress := string(c.Request.Header.Peek("X-Real-IP"))
	remoteIpAddress := string(c.Request.Header.Peek("X-Real-IP"))

	if remoteAddress == "" {
		remoteAddress = c.RemoteAddr().String()
		remoteIpAddress = c.RemoteIP().String()
	}

	remoteHost := string(c.Request.Header.Peek("Host"))
	if remoteHost == "" {
		remoteHost = string(c.Request.Host())
	}

	return models.RequestLogging{string(
		c.Method()),
		string(c.Request.RequestURI()),
		string(c.Request.Body()),
		remoteHost,
		string(c.UserAgent()),
		remoteAddress,
		remoteIpAddress,
		string(headers),
		string(queryArgs),
		currentDate,
		string(c.Response.Body()),
		elapsedTime,
		c.Response.StatusCode(),
		serviceName,
		indexName}
}

func PublishLog(reqLogging *models.RequestLogging) {
	defer utils.PreventCrash()

	LoggingType[config.GApiConfiguration.Logs.Type].(func(*models.RequestLogging))(reqLogging)
	return
}

func PublishRabbit(reqLogging *models.RequestLogging) {
	reqLoggingJson, _ := json.Marshal(reqLogging)

	/* if RabbitConnection == nil{
		RabbitConnection = rabbit.ConnectToRabbit()
	} */

	// RabbitChannel := rabbit.CreateChannel(rabbit.RabbitConnection)
	/*
		LogQueue, err := RabbitChannel.QueueDeclare(
			rabbit.Queue(), // name
			true,           // durable
			false,          // delete when unused
			false,          // exclusive
			false,          	// no-wait
			nil,            // arguments
		)
		rabbit.FailOnError(err, "Failed to declare queue")
	*/
	/* err := rabbit.RabbitChannelGlobal.Publish(
	"",            // exchange
	rabbit.Queue(), // routing key
	false,         // mandatory
	false,         	// immediate
	amqp.Publishing{
		ContentType: "application/json",
		Body:        reqLoggingJson,
	}) */

	rabbit.PublishToRabbitMQ(reqLoggingJson)

	// rabbit.FailOnError(err, "Failed to publish message")
}
