package http

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"strconv"
	"errors"
)

// ParsePageParam parses page query parameter
func ParsePageParam(c *routing.Context) (int, error){
	page := 1
	if c.QueryArgs().Has("page") {
		page, err := strconv.Atoi(string(c.QueryArgs().Peek("page")))

		if err != nil {
			return -1, errors.New("Invalid page provided")
		}
	}
	return page, nil
}