package system1

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	ProductRouter
}

var exaProductRouter = api.ApiGroupApp.System1ApiGroup.ProductsApi
