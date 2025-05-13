package products

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	service "github.com/flipped-aurora/gin-vue-admin/server/service/shopping"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

type ProductsApi struct{}

// GetProductSku 获取商品 SKU 信息
// @Tags 商品管理
// @Security ApiKeyAuth
// @Summary 根据 sku_id 和 product_code 获取商品 SKU 详情
// @Param sku_id query string false "SKU ID"
// @Param product_code query string true "商品编码"
// @Success 200 {object} response.Response "返回商品 SKU 信息"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/targetSku_id [get]  // 更新这里的路径
func (s *ProductsApi) GetProductSku(c *gin.Context) {
	skuId := c.Query("sku_id")
	productCode := c.Query("product_code")
	sku, err := service.GetProductInfoList(skuId, productCode)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(sku, c)
}

// GetProductReview 获取商品的评论信息
// @Tags 商品管理
// @Summary 根据 product_code 获取商品评论列表
// @Description 获取指定商品的评论信息，支持分页、排序、星级筛选
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param product_code query string true "商品编码（必填）"
// @Param page query int false "页码（默认 1）"
// @Param limit query int false "每页数量（默认 10，最大100）"
// @Param sort query string false "排序方式（newest, oldest, highest_rating, lowest_rating, most_helpful）"
// @Param rating query int false "星级筛选（1～5）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/reviews [get]
func (s *ProductsApi) GetProductReview(c *gin.Context) {
	productCode := c.Query("product_code")
	pageurl := c.DefaultQuery("page", "1") //从 URL 的 query string 里取参数如果没有这个参数，就用你给的默认值  “，”后面的就是你定好的默认值
	limiturl := c.DefaultQuery("limit", "10")
	sorturl := c.DefaultQuery("sort", "newest")
	ratingurl := c.DefaultQuery("rating", "0")

	// if productCode == "" {
	// 	response.FailWithMessage("不正な商品識別子です", c)
	// 	return
	// }
	page, err := strconv.Atoi(pageurl)
	// strconv 是 Go 标准库里的一个包，专门用来做字符串和数字之间的转换。
	// Atoi 是它的一个函数，全称是 "ASCII to Integer"，意思是把字符串变成整数。
	// 你传进去一个字符串，比如 "3"，它就会返回数字 3。
	// 如果传进去 "abc" 这种乱七八糟的，转换失败，它就会返回一个 error。
	// 所以 strconv.Atoi() 是：字符串 ➔ 整数（int） 的转换。
	// page 是转换成功后的整数值。
	// err 是转换过程中出现的错误（如果有的话）。
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			//http.StatusBadRequest  HTTP状态码400，意思是前端请求错了 400 表示是客户端的问题（比如传错参数）
			"error": "pageパラメータは数値で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	if page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pageパラメータは1以上で指定してください",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	limit, err := strconv.Atoi(limiturl)
	if err != nil || limit < 1 || limit >= 101 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limitパラメータは1から100の間で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	if sorturl != "newest" && sorturl != "oldest" && sorturl != "highest_rating" && sorturl != "lowest_rating" && sorturl != "most_helpful" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不正なsortパラメータです。('newest', 'oldest', 'highest_rating', 'lowest_rating', 'most_helpful' のいずれかを指定)",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	rating, err := strconv.Atoi(ratingurl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ratingパラメータは1から5の間で指定してください",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	if rating > 5 || rating < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ratingパラメータは1から5の間で指定してください",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	reviewService := &service.Getreview{}
	neededReview, err := reviewService.GetProductReviewInfo(productCode, page, limit, sorturl, rating)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		//response.FailWithMessage("不正な商品識別子です", c)
		return
	}

	response.OkWithData(neededReview, c)
}

// GetProductQueation 获取商品的问答信息
// @Tags 商品管理
// @Summary 根据 product_code 获取商品问答列表
// @Description 获取指定商品的问答信息，支持分页、排序
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param product_code query string true "商品编码（必填）"
// @Param page query int false "页码（默认 1）"
// @Param limit query int false "每页数量（默认 10，最大100）"
// @Param sort query string false "排序方式（newest, oldest, most_helpful）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/questions [get]
func (s *ProductsApi) GetProductQueation(c *gin.Context) {
	productCode := c.Query("product_code")
	pageurl := c.DefaultQuery("page", "1")
	limiturl := c.DefaultQuery("limit", "10")
	sorturl := c.DefaultQuery("sort", "newest")
	page, err := strconv.Atoi(pageurl)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pageパラメータは数値で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	if page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pageパラメータは1以上で指定してください",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	limit, err := strconv.Atoi(limiturl)
	if err != nil || limit < 1 || limit >= 101 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limitパラメータは1から100の間で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	if sorturl != "newest" && sorturl != "oldest" && sorturl != "most_helpful" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不正なsortパラメータです。('newest', 'oldest', 'most_helpful' のいずれかを指定)",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	qesService := &service.Getquestion{}
	neededQuestion, err := qesService.GetProductQueationInfo(productCode, page, limit, sorturl)

	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		//response.FailWithMessage("不正な商品識別子です", c)
		return
	}
	response.OkWithData(neededQuestion, c)
}

// GetProductImages 获取商品画像
// @Tags 商品管理
// @Summary 根据 sku_id 获取商品画像
// @Description 获取指定商品的画像
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param sku_id query string true "商品编码（必填）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/images [get]
func (s *ProductsApi) GetProductsImages(c *gin.Context) {
	skuId := c.Query("sku_id")
	imageService := &service.GetImages{}
	sku1, err := imageService.GetProductsImagesInfo(skuId)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(sku1, c)
}

// Addfavorite 添加收藏列表
// @Tags 商品管理
// @Summary 根据 sku_id 添加收藏
// @Description 根据 sku_id 添加收藏
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param sku_id query string true "商品编码（必填）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/favorites/add [post]
func (s *ProductsApi) Addfavorite(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		// 如果 uid 为0，说明没拿到，返回未登录
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	skuId := c.Query("sku_id")
	addService := &service.FavoriteList{}
	add, err := addService.AddfavoriteList(fmt.Sprintf("%d", uid), skuId)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(add, c)
}

// Deletefavorite 删除收藏列表
// @Tags 商品管理
// @Summary 根据 sku_id 删除收藏
// @Description 根据 sku_id 删除收藏
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param sku_id query string true "商品编码（必填）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/favorites/delete [delete]
func (s *ProductsApi) Deletefavorite(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		// 如果 uid 为0，说明没拿到，返回未登录
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	skuId := c.Query("sku_id")
	delService := &service.FavoriteList{}
	del, err := delService.DELETEfavoriteList(fmt.Sprintf("%d", uid), skuId)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(del, c)
}

// GetFavoriteList 获取用户收藏列表
// @Tags 商品管理
// @Summary 获取用户收藏列表
// @Description 获取指定用户收藏列表
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param page query int false "页码（默认 1）"
// @Param limit query int false "每页数量（默认 10，最大100）"
// @Param sort query string false "排序方式（newest, oldest）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/favorites/skus [get]
func (s *ProductsApi) GetFavoriteList(c *gin.Context) {
	uid := utils.GetUserID(c)
	pageurl := c.DefaultQuery("page", "1")
	limiturl := c.DefaultQuery("limit", "10")
	sorturl := c.DefaultQuery("sort", "newest")

	if uid == 0 {
		response.FailWithDetailed(nil, "ユーザーIDが指定されていません。", c)
		return
	}
	page, err := strconv.Atoi(pageurl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pageパラメータは数値で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	if page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pageパラメータは1以上で指定してください",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	limit, err := strconv.Atoi(limiturl)
	if err != nil || limit < 1 || limit >= 101 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limitパラメータは1から100の間で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	if sorturl != "newest" && sorturl != "oldest" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不正なsortパラメータです。('newest', 'oldest',  のいずれかを指定)",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	likeService := &service.FavoriteList{}
	likeList, err := likeService.GetFavoriteListInfo(fmt.Sprintf("%d", uid), page, limit, sorturl)

	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		//response.FailWithMessage("不正な商品識別子です", c)
		return
	}
	response.OkWithData(likeList, c)
}

// GetRelatedList 获取关联商品
// @Tags 商品管理
// @Summary 获取关联商品
// @Description 根据商品code获得关联商品
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param product_code query string true "商品编码（必填）"
// @Param limit query int false "每页数量（默认 5，最大100）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/related [get]
func (s *ProductsApi) GetRelatedList(c *gin.Context) {
	productCode := c.Query("product_code")
	limiturl := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limiturl)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limitパラメータは1以上の数値で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}

	relatedService := &service.GetRelatedList{}
	relatedList, err := relatedService.GetRelatedListInfo(productCode, limit)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(relatedList, c)
}

// GetCoordinateList 获取工作人员的推荐搭配
// @Tags 商品管理
// @Summary 获取相关搭配商品
// @Description 根据商品code获得推荐搭配
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param product_code query string true "商品编码（必填）"
// @Param limit query int false "每页数量（默认 4，最大100）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/products/coordinates [get]
func (s *ProductsApi) GetCoordinateList(c *gin.Context) {
	productCode := c.Query("product_code")
	limiturl := c.DefaultQuery("limit", "4")

	limit, err := strconv.Atoi(limiturl)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limitパラメータは1以上の数値で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	coordinateService := &service.GetCoordinateList{}
	coordinateList, err := coordinateService.GetCoordinateListInfo(productCode, limit)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(coordinateList, c)
}

// AddCheckList 添加浏览历史
// @Tags 用户浏览记录
// @Summary 根据 sku_id 添加浏览历史
// @Description 根据 sku_id 添加浏览历史
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param sku_id query string true "商品编码（必填）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/history/add [post]
func (s *ProductsApi) AddHistoryList(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		// 如果 uid 为0，说明没拿到，返回未登录
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	skuId := c.Query("sku_id")
	AddCheckService := &service.GetHistoryList{}
	addCheck, err := AddCheckService.AddHistoryList(fmt.Sprintf("%d", uid), skuId)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(addCheck, c)
}

// DeleteCheckList 删除浏览历史
// @Tags 用户浏览记录
// @Summary 根据 sku_id 删除浏览历史
// @Description 根据 sku_id 删除浏览历史
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param sku_id query string true "商品编码（必填）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/history/delete [delete]
func (s *ProductsApi) DeleteCheck(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		// 如果 uid 为0，说明没拿到，返回未登录
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	skuId := c.Query("sku_id")
	delCheckService := &service.GetHistoryList{}
	delCheck, err := delCheckService.DELETECheckList(fmt.Sprintf("%d", uid), skuId)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(delCheck, c)
}

// GetCheckList 获取用户浏览记录
// @Tags 用户浏览记录
// @Summary 获取用户浏览记录
// @Description 获得用户浏览记录
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/history [get]
func (s *ProductsApi) GetCheckList(c *gin.Context) {
	uid := utils.GetUserID(c)
	pageurl := c.DefaultQuery("page", "1")
	limiturl := c.DefaultQuery("limit", "10")
	page, err := strconv.Atoi(pageurl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pageパラメータは数値で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	if page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "pageパラメータは1以上で指定してください",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	limit, err := strconv.Atoi(limiturl)
	if err != nil || limit < 1 || limit >= 101 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limitパラメータは1から100の間で指定してください。",
			"code":  "INVALID_PARAMETER",
		})
		return
	}
	HistoryService := &service.GetHistoryList{}
	historyList, err := HistoryService.GetCheckListInfo(fmt.Sprintf("%d", uid), page, limit)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	if historyList == nil {
		response.FailWithMessage("データが取得できませんでした。", c)
		return
	}
	response.OkWithData(historyList, c)
}

// Addcart 添加购物车
// @Tags 用户购物车
// @Summary 输入 sku_id和数量来添加购物车
// @Description 输入 sku_id和数量来添加购物车
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param sku_id query string true "商品编码（必填）"
// @Param quantity query string false "商品数量，不填则为1"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/cart/items/add [post]
func (s *ProductsApi) AddCartList(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		// 如果 uid 为0，说明没拿到，返回未登录
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	skuId := c.Query("sku_id")
	quantityUrl := c.DefaultQuery("quantity", "1")
	quantity, err := strconv.Atoi(quantityUrl)
	if err != nil || quantity <= 0 || quantity > 999 {
		response.FailWithMessage("数量は1以上、999以下で設定してください", c)
	}
	addCartService := &service.GetCart{}
	addCart, err := addCartService.AddCartListInfo(fmt.Sprintf("%d", uid), skuId, quantity)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(addCart, c)
}

// delcart 删除购物车
// @Tags 用户购物车
// @Summary 输入 sku_id来删除购物车
// @Description 输入 sku_id来删除购物车
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param sku_id query string true "商品编码（必填）"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/cart/items/del [delete]
func (s *ProductsApi) DelCartList(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		// 如果 uid 为0，说明没拿到，返回未登录
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	skuId := c.Query("sku_id")
	delCartService := &service.GetCart{}
	delCart, err := delCartService.DelCartListInfo(fmt.Sprintf("%d", uid), skuId)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(delCart, c)
}

// CartQuantityChange 更改购物车商品数量
// @Tags 用户购物车
// @Summary 输入 sku_id和数量来更改商品数量
// @Description 输入 sku_id和数量来更改商品数量
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param sku_id query string true "商品编码（必填）"
// @Param quantity query string true "订单数量 (必填)必须是1以上"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/cart/items/change [put]
func (s *ProductsApi) CartQuantityChange(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		// 如果 uid 为0，说明没拿到，返回未登录
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	skuId := c.Query("sku_id")
	quantityUrl := c.Query("quantity")
	quantity, err := strconv.Atoi(quantityUrl)
	if err != nil || quantity <= 0 || quantity > 999 {
		response.FailWithMessage("数量は1以上、999以下で設定してください", c)
	}
	changeCartService := &service.GetCart{}
	delCart, err := changeCartService.CartQuantityChangeInfo(fmt.Sprintf("%d", uid), skuId, quantity)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(delCart, c)
}

// GetCartList 获取用户购物车列表
// @Tags 用户购物车
// @Summary 获取用户购物车列表
// @Description 获得用户购物车列表
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/cart [get]
func (s *ProductsApi) GetCartList(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	getCartService := &service.GetCart{}
	getCart, err := getCartService.GetCartListInfo(fmt.Sprintf("%d", uid))
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(getCart, c)
}

// GetAddressList 获取用户配送地址列表
// @Tags 用户配送地址
// @Summary 用户配送地址列表
// @Description 用户配送地址列表
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/shipping-addresses [get]
func (s *ProductsApi) GetAddressList(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	getAdressService := &service.Address{}
	getAdress, err := getAdressService.GetAddressListInfo(fmt.Sprintf("%d", uid))
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(getAdress, c)
}

// AddAddressList 追加用户配送地址
// @Tags 用户配送地址
// @Summary 追加用户配送地址
// @Description 追加用户配送地址
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param postal_code query string true "邮政编码"
// @Param prefecture query string true "都道府県"
// @Param city query string true "市区町村"
// @Param address_line1 query string true "丁目・番地"
// @Param address_line2 query string true "建物名・部屋番号"
// @Param recipient_name query string true "氏名"
// @Param phone_number query string true "電話番号"
// @Param is_default query string true "是否设为默认地址1为是，2为否"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/shipping-addresses/add [post]
func (s *ProductsApi) AddAddressList(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	PostalCode := c.Query("postal_code")
	Prefecture := c.Query("prefecture")
	City := c.Query("city")
	AddressLine1 := c.Query("address_line1")
	AddressLine2 := c.Query("address_line2")
	RecipientName := c.Query("recipient_name")
	PhoneNumber := c.Query("phone_number")
	IsDefault := c.Query("is_default")
	addAdressService := &service.Address{}
	isDefaultInt, err := strconv.Atoi(IsDefault)
	if err != nil {
		response.FailWithMessage("is_defaultパラメータは数値で指定してください。", c)
		return
	}
	addAdress, err := addAdressService.AddAddressListInfo(fmt.Sprintf("%d", uid), PostalCode, Prefecture, City, AddressLine1, AddressLine2, RecipientName, PhoneNumber, isDefaultInt)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(addAdress, c)
}

// DelAddressList 删除用户配送地址
// @Tags 用户配送地址
// @Summary 删除用户配送地址
// @Description 删除用户配送地址
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param address_id query string true "address_id"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/shipping-addresses/del [delete]
func (s *ProductsApi) DelAddressList(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	addressId := c.Query("address_id")
	delAdressService := &service.Address{}
	delAdress, err := delAdressService.DelAddressListInfo(fmt.Sprintf("%d", uid), addressId)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(delAdress, c)
}

// AddAddressList 编辑用户配送地址
// @Tags 用户配送地址
// @Summary 编辑用户配送地址
// @Description 编辑用户配送地址
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param address_id query string true "地址编号"
// @Param postal_code query string false "邮政编码"
// @Param prefecture query string false "都道府県"
// @Param city query string false "市区町村"
// @Param address_line1 query string false "丁目・番地"
// @Param address_line2 query string false "建物名・部屋番号"
// @Param recipient_name query string false "氏名"
// @Param phone_number query string false "電話番号"
// @Param is_default query string false "是否设为默认地址1为是，2为否"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/shipping-addresses/change [put]
func (s *ProductsApi) ChangeAddressList(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	addressId := c.Query("address_id")
	PostalCode := c.Query("postal_code")
	Prefecture := c.Query("prefecture")
	City := c.Query("city")
	AddressLine1 := c.Query("address_line1")
	AddressLine2 := c.Query("address_line2")
	RecipientName := c.Query("recipient_name")
	PhoneNumber := c.Query("phone_number")
	IsDefault := c.Query("is_default")
	changeAdressService := &service.Address{}
	isDefaultInt := 0 // 先给默认值0
	var err error
	if IsDefault != "" {
		isDefaultInt, err = strconv.Atoi(IsDefault)
		if err != nil {
			response.FailWithMessage("is_defaultパラメータは数値で指定してください。", c)
			return
		}
	}

	changeAdress, err := changeAdressService.ChangeAddressInfo(fmt.Sprintf("%d", uid), addressId, PostalCode, Prefecture, City, AddressLine1, AddressLine2, RecipientName, PhoneNumber, isDefaultInt)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(changeAdress, c)
}

// GetAddressList 获取支付方式
// @Tags 支付
// @Summary 支付方式
// @Description 支付方式
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/payments/methods [get]
func (s *ProductsApi) GetPaymentList(c *gin.Context) {
	getPaymentService := &service.Payment{}
	getPayment, err := getPaymentService.GetPaymentListInfo()
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(getPayment, c)
}

// GetCouponPoint 获取可利用优惠券和积分
// @Tags 支付
// @Summary 获取可利用优惠券和积分
// @Description 获取可利用优惠券和积分
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/orders/checkout/info [get]
func (s *ProductsApi) GetCouponPoint(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	getCouponPointService := &service.Payment{}
	getCouponPoint, err := getCouponPointService.GetCouponPointListInfo(fmt.Sprintf("%d", uid))
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(getCouponPoint, c)
}

// GetCouponPoint 使用优惠券
// @Tags 支付
// @Summary 使用优惠券
// @Description 使用优惠券
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param coupon_code query string true "优惠券code"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/orders/checkout/apply-coupon [post]
func (s *ProductsApi) UseCoupon(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	couponCode := c.Query("coupon_code")
	useCouponService := &service.Payment{}
	useCoupon, err := useCouponService.UseCouponInfo(fmt.Sprintf("%d", uid), couponCode)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(useCoupon, c)
}

// GetCouponPoint 撤销优惠券
// @Tags 支付
// @Summary 撤销优惠券
// @Description 撤销优惠券
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param coupon_code query string true "优惠券code"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/orders/checkout/remove-coupon [delete]
func (s *ProductsApi) RemoveCoupon(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	couponCode := c.Query("coupon_code")
	removeCouponService := &service.Payment{}
	removeCoupon, err := removeCouponService.RemoveCouponInfo(fmt.Sprintf("%d", uid), couponCode)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(removeCoupon, c)
}

// GetCouponPoint 使用积分
// @Tags 支付
// @Summary 使用积分
// @Description 使用积分
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param points query string true "输入要消耗的积分"
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/orders/checkout/use-points [post]
func (s *ProductsApi) UsePoints(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	// point := c.Query("available_points")
	pointStr := c.Query("points")
	point, err := strconv.Atoi(pointStr)
	if err != nil || point <= 0 {
		response.FailWithDetailed(nil, "参数 available_points 无效", c)
		return
	}
	usePointsService := &service.Payment{}
	usePoints, err := usePointsService.UsePointInfo(fmt.Sprintf("%d", uid), point)
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(usePoints, c)
}

// GetCouponPoint 撤销使用的积分
// @Tags 支付
// @Summary 撤销使用的积分
// @Description 撤销使用的积分
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "返回评论列表"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/orders/checkout/remove-points [delete]
func (s *ProductsApi) RemovePoints(c *gin.Context) {
	uid := utils.GetUserID(c)
	if uid == 0 {
		response.FailWithDetailed(nil, "用户未登录或token无效", c)
		return
	}
	removePointsService := &service.Payment{}
	removePoints, err := removePointsService.RemovePointInfo(fmt.Sprintf("%d", uid))
	if err != nil {
		response.FailWithDetailed(nil, err.Error(), c)
		return
	}
	response.OkWithData(removePoints, c)
}
