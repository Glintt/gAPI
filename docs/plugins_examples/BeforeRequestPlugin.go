package main

import (
	"errors"
	"fmt"

	routing "github.com/qiangxue/fasthttp-routing"
)

type beforeRequestEntryPlugin string

func (g beforeRequestEntryPlugin) Call(c *routing.Context) error {
	fmt.Println("PLUGIN ---- BeforeRequestPlugin")

	c.Response.SetBody([]byte(`{"error":true, "msg": "Not authorized to access resource."}`))
	c.Response.SetStatusCode(405)
	c.Response.Header.SetContentType("application/json")
	c.Abort()

	return errors.New("Return an error which will stop the request")
}

var BeforeRequestEntryPlugin beforeRequestEntryPlugin
