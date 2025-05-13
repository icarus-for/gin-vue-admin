package system1

import (
	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (s *ProductRouter) InitProductRouter(Router *gin.RouterGroup) {
	ProductRouter := Router.Group("/api/v1")
	ProductRouter.GET("/products/targetSku_id", exaProductRouter.GetProductSku)
	ProductRouter.GET("/products/reviews", exaProductRouter.GetProductReview)
	ProductRouter.GET("/products/questions", exaProductRouter.GetProductQueation)
	ProductRouter.GET("/products/images", exaProductRouter.GetProductsImages)
	ProductRouter.GET("/products/favorites/skus", exaProductRouter.GetFavoriteList)
	ProductRouter.POST("/products/favorites/add", exaProductRouter.Addfavorite)
	ProductRouter.DELETE("/products/favorites/delete", exaProductRouter.Deletefavorite)
	ProductRouter.GET("/products/related", exaProductRouter.GetRelatedList)
	ProductRouter.GET("/products/coordinates", exaProductRouter.GetCoordinateList)

	ProductRouter.GET("/history", exaProductRouter.GetCheckList)
	ProductRouter.DELETE("/history/delete", exaProductRouter.DeleteCheck)
	ProductRouter.POST("/history/add", exaProductRouter.AddHistoryList)

	ProductRouter.GET("/cart", exaProductRouter.GetCartList)
	ProductRouter.POST("/cart/items/add", exaProductRouter.AddCartList)
	ProductRouter.DELETE("/cart/items/del", exaProductRouter.DelCartList)
	ProductRouter.PUT("/cart/items/change", exaProductRouter.CartQuantityChange)

	ProductRouter.GET("/shipping-addresses", exaProductRouter.GetAddressList)
	ProductRouter.POST("/shipping-addresses/add", exaProductRouter.AddAddressList)
	ProductRouter.DELETE("/shipping-addresses/del", exaProductRouter.DelAddressList)
	ProductRouter.PUT("/shipping-addresses/change", exaProductRouter.ChangeAddressList)

	ProductRouter.GET("/payments/methods", exaProductRouter.GetPaymentList)
	ProductRouter.GET("/orders/checkout/info", exaProductRouter.GetCouponPoint)
	ProductRouter.POST("/orders/checkout/apply-coupon", exaProductRouter.UseCoupon)
	ProductRouter.DELETE("/orders/checkout/remove-coupon", exaProductRouter.RemoveCoupon)
	ProductRouter.POST("/orders/checkout/use-points", exaProductRouter.UsePoints)
	ProductRouter.DELETE("/orders/checkout/remove-points", exaProductRouter.RemovePoints)

}
