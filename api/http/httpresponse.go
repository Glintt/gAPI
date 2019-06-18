package http

import (
	"github.com/Glintt/gAPI/api/utils"

	routing "github.com/qiangxue/fasthttp-routing"
)

type ResponseInfo struct {
	StatusCode  int
	ContentType []byte
	Body        []byte
}

// Response creates the API response structure 
func Response(c *routing.Context, response string, statuscode int, service string, contentType string) {
	utils.LogMessage("RESPONSE ==> "+response, utils.DebugLogType)

	c.Response.SetBody([]byte(response))
	c.Response.Header.SetContentType(contentType)
	c.Response.Header.Set("service", service)
	c.Response.Header.SetStatusCode(statuscode)
}

// Error creates and sends error response
func Error(c *routing.Context, err string, statusCode int, service string) error {
	Response(c, `{"error": true, "msg": "` + err + `"}`, statusCode, service, "application/json")
	return nil
}

// Created Creates and sends created response
func Created(c *routing.Context, msg string, service string) error {
	Response(c, `{"error": false, "msg": "`+msg+`"}`, 201, service, "application/json")
	return nil
}

// OkFormated creates and send created response
func OkFormated(c *routing.Context, msg string, service string) error {
	Response(c, `{"error": false, "msg": "`+msg+`"}`, 200, service, "application/json")
	return nil
}

// Ok creates and send created response
func Ok(c *routing.Context, msg string, service string) error {
	Response(c, msg, 200, service, "application/json")
	return nil
}

// Deleted creates and send deleted response
func Deleted(c *routing.Context, msg string, service string) error {
	Response(c, `{"error": false, "msg": "`+msg+`"}`, 200, service, "application/json")
	return nil
}

// NotFound creates and send not found response
func NotFound(c *routing.Context, msg string, service string) error {
	Response(c, `{"error": true, "msg": "`+msg+`"}`, 404, service, "application/json")
	return nil
}
