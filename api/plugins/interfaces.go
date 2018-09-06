package plugins

import routing "github.com/qiangxue/fasthttp-routing"

type BeforeRequestEntryPlugin interface {
	Call(c *routing.Context) error
}
