package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	system1 "github.com/flipped-aurora/gin-vue-admin/server/router/products"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	System1 system1.RouterGroup
}
