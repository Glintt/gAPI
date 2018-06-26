package logs

import (
	"gAPIManagement/api/config"
	"gAPIManagement/api/rabbit"
	"gAPIManagement/api/utils"
	"encoding/json"

	// "github.com/streadway/amqp"

	"github.com/valyala/fasthttp"
)
/* 
var RabbitConnection *amqp.Connection
var RabbitConnection *amqp.Connection
var RabbitChannel *amqp.Channel
var LogQueue amqp.Queue */

type RequestLogging struct {
	Method      string
	Uri         string
	RequestBody string
	Host        string
	UserAgent   string
	RemoteAddr  string
	RemoteIp    string
	Headers     string
	QueryArgs   string
	DateTime    string
	Response    string
	ElapsedTime int64
	StatusCode  int
	ServiceName string
}

var LoggingType = map[string]interface{}{
	"Rabbit": PublishRabbit,
	"Elastic": PublishElastic}

func NewRequestLogging(c *fasthttp.RequestCtx, queryArgs []byte, headers []byte, currentDate string, elapsedTime int64, serviceName string) RequestLogging {
	return RequestLogging{string(
		c.Method()),
		string(c.Request.RequestURI()),
		string(c.Request.Body()),
		string(c.Request.Host()),
		string(c.UserAgent()),
		c.RemoteAddr().String(),
		c.RemoteIP().String(),
		string(headers),
		string(queryArgs),
		currentDate,
		string(c.Response.Body()),
		elapsedTime,
		c.Response.StatusCode(),
		serviceName}
}

func (reqLogging *RequestLogging) Save() {	
	go PublishLog(reqLogging)
}

func PublishLog(reqLogging *RequestLogging) {
	defer utils.PreventCrash()

	LoggingType[config.GApiConfiguration.Logs.Type].(func(*RequestLogging))(reqLogging)
}

func PublishElastic(reqLogging *RequestLogging) {
	utils.LogMessage("ELASTIC PUBLISH", utils.DebugLogType)
	currentDate := utils.CurrentDate()
	logsURL := config.ELASTICSEARCH_URL + "/gapi-logs/request-logs-" + currentDate

	reqLoggingJson, _ := json.Marshal(reqLogging)

	request := fasthttp.AcquireRequest()

	request.SetRequestURI(logsURL)
	request.Header.SetMethod("POST")

	request.Header.SetContentType("application/json")
	request.SetBody(reqLoggingJson)
	client := fasthttp.Client{}

	resp := fasthttp.AcquireResponse()
	err := client.Do(request, resp)

	if err != nil {
		utils.LogMessage("ELASTIC PUBLISH - error:" + err.Error(), utils.DebugLogType)
		
		resp.SetStatusCode(400)
	}
}

func PublishRabbit(reqLogging *RequestLogging) {
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