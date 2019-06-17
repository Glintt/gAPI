package main

import (
	"errors"

	"github.com/Glintt/gAPI/api/utils"
	routing "github.com/qiangxue/fasthttp-routing"
)

type beforeRequestEntryPlugin string

func (g beforeRequestEntryPlugin) Call(c *routing.Context) error {
	utils.LogMessage(err.Error(), utils.DebugLogType)

	c.Response.SetBody([]byte(`{"error":true, "msg": "Not authorized to access resource."}`))
	c.Response.SetStatusCode(405)
	c.Response.Header.SetContentType("application/json")
	c.Abort()

	return errors.New("Return an error which will stop the request")
}

var BeforeRequestEntryPlugin beforeRequestEntryPlugin
