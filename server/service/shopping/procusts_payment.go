package service

import (
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
)

// type tempCAP struct {
// 	CouponID        uint64  `json:"coupon_id"`
// 	CouponCode      string  `json:"coupon_code"`
// 	Name            string  `json:"name"`
// 	Description     *string `json:"description,omitempty"`
// 	DiscountText    string  `json:"discount_text,omitempty"`
// 	AvailablePoints int     `json:"available_points"`
// }

type Payment struct{}

func (p *Payment) GetPaymentListInfo() ([]dto.PaymentMethodInfo, error) {
	var paymentList []dto.PaymentMethodInfo
	err := global.GVA_DB.Table("payment_methods").
		Select(`
			id AS method_id,
			method_code,
			name,
			description
			`).
		Where("is_active = ?", "1").
		Scan(&paymentList).Error
	if err != nil {
		return nil, err
	}
	return paymentList, err
}

func (p *Payment) GetCouponPointListInfo(userId string) (*dto.CheckoutInfoResponseV1_1, error) {
	var couponList []dto.AvailableCouponInfo
	err := global.GVA_DB.Table("user_coupons").
		Select(`
			user_coupons.coupon_id,
			coupons.coupon_code,
			coupons.name,
			coupons.description,
			CONCAT(
			 	'减少',discount_value,' ',
			CASE
				WHEN discount_type = 'percentage ' THEN '%'
				WHEN discount_type = 'fixed ' THEN '円'
				ELSE discount_type
			END,
				'最低購入金額は ',min_purchase_amount, '円',
				'利用上限額 ',max_discount_amount, '円'
				)AS discount_text
			`).
		Joins("LEFT JOIN coupons ON user_coupons.coupon_id = coupons.id").
		Where("user_coupons.user_id= ? AND user_coupons.is_used = 0 AND NOW() BETWEEN coupons.start_date AND coupons.end_date", userId).
		Scan(&couponList).Error
	if err != nil {
		return nil, err
	}

	//查出userid对应的购物车价格
	type tempAmount struct {
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
	}
	var temp []tempAmount
	err = global.GVA_DB.Table("user_cart_items").
		Distinct().
		Select(`
			user_cart_items.quantity,
               COALESCE(
            (
                SELECT p.price
                FROM prices p
                WHERE p.sku_id = user_cart_items.sku_id
                AND p.price_type_id IN (2,3)
                AND (p.start_date IS NULL OR p.start_date <= NOW())
                AND (p.end_date IS NULL OR p.end_date >= NOW())
                LIMIT 1
            ),
            (
                SELECT p2.price
                FROM prices p2
                WHERE p2.sku_id = user_cart_items.sku_id
                AND p2.price_type_id = 1
                LIMIT 1
            )
        ) AS price 
			`).
		Joins("LEFT JOIN prices ON user_cart_items.sku_id = prices.sku_id").
		Joins("LEFT JOIN price_types ON prices.price_type_id = price_types.id").
		Where("user_cart_items.user_id=?", userId).
		Scan(&temp).Error
	if err != nil {
		return nil, err
	}
	var TotalItemsCount int
	var TotalAmount float64

	for _, v := range temp {
		TotalItemsCount += v.Quantity
		TotalAmount += v.Price * float64(v.Quantity)
	}
	//更新到checkout_sessions表里面
	err = global.GVA_DB.Table("checkout_sessions").
		Where("user_id = ?", userId).
		Update("cart_subtotal", TotalAmount).Error
	if err != nil {
		return nil, err
	}

	var userPoint int
	err = global.GVA_DB.Table("user_points").
		Select("available_points").
		Where("user_id=?", userId).
		Scan(&userPoint).Error
	if err != nil {
		return nil, err
	}

	// type tempStruct struct {
	// 	CartSubtotalAmount string  `json:"cart_subtotal"`
	// 	CouponID           uint64  `json:"coupon_id"`
	// 	CouponCode         string  `json:"coupon_code"`
	// 	Name               string  `json:"name"`
	// 	DiscountAmount     float64 `json:"discount_amount"` // この注文での実際の割引額 (計算用)
	// 	//FormattedDiscountAmount     string  `json:"formatted_discount_amount"` // 表示用
	// 	PointsDiscountAmountFormatted string `json:"points_discount_amount_formatted"` // ポイント割引額 (表示用)
	// 	ShippingFeeFormatted          string `json:"shipping_fee_formatted"`           // 送料 (表示用、別途計算の場合あり)
	// 	TotalAmountFormatted          string `json:"total_amount_formatted"`
	// }

	// //积分和券都使用，去查被写入新数据的checkout_sessions的这个表
	// var currentCheckoutState []tempStruct
	// err = global.GVA_DB.Table("checkout_sessions").
	// 	Where("user_id=?", userId).
	// 	Scan(&currentCheckoutState).Error
	// if err != nil {
	// 	return nil, fmt.Errorf("查找数据失败")
	// }

	cAndP := &dto.CheckoutInfoResponseV1_1{
		AvailableCoupons: couponList,
		UserPoints: &dto.UserPointInfo{
			AvailablePoints: userPoint,
		},
	}
	return cAndP, err
}

func (p *Payment) UseCouponInfo(userId string, couponCode string) (*dto.CheckoutInfoResponseV1_1, error) {
	//先找有没有这个优惠券
	var couponFind dto.AvailableCouponInfo
	err := global.GVA_DB.Table("coupons").
		Select(`
				id AS coupon_id,
				coupon_code,
				name,
				discount_value,
				discount_type,
				min_purchase_amount,
				max_discount_amount
				`).
		Where("coupon_code = ? AND NOW() BETWEEN start_date AND end_date", couponCode).
		Scan(&couponFind).Error
	if err != nil {
		return nil, fmt.Errorf("优惠券不存在或已过期")
	}
	//在查用户有没有这个优惠券
	var userCoupon int64
	err = global.GVA_DB.Table("user_coupons").
		Select("coupon_id").
		Where("user_id = ? AND coupon_id = ? AND is_used = 0", userId, couponFind.CouponID).
		Count(&userCoupon).Error
	if err != nil || userCoupon == 0 {
		return nil, fmt.Errorf("此优惠券不可用或已使用")
	}

	//查优惠前购物车总金额和被使用的积分
	var cartInformation struct {
		CartTotalAmount float64 `gorm:"column:cart_subtotal"`
		UsedPoints      int     `gorm:"column:used_points"`
		ShippingFee     float64 `gorm:"column:shipping_fee"`
	}
	err = global.GVA_DB.Table("checkout_sessions").
		Select("cart_subtotal,used_points,shipping_fee").
		Where("user_id=?", userId).
		Scan(&cartInformation).Error
	if err != nil {
		return nil, fmt.Errorf("获取购物车信息失败")
	}
	//看看是否达到最小使用金额
	if cartInformation.CartTotalAmount < couponFind.MinPurchaseAmount {
		return nil, fmt.Errorf("购物金额未达到优惠券使用条件")
	}
	//要先查使用条件，否则即使使用不了，也会被扣除
	//如果用了 就给is_used变为1，日期也要更新
	err = global.GVA_DB.Table("user_coupons").
		Where("user_id = ? AND coupon_id = ?", userId, couponFind.CouponID).
		Updates(map[string]interface{}{
			"is_used": 1,
			"used_at": time.Now(),
		}).Error
	if err != nil {
		return nil, fmt.Errorf("该优惠券已经被使用")
	}

	//计算优惠之后的金额
	discountAmount := p.calculateDiscount(couponFind.DiscountType, couponFind.DiscountValue, couponFind.MaxDiscountAmount, cartInformation.CartTotalAmount)
	totalAmount := cartInformation.CartTotalAmount + cartInformation.ShippingFee - float64(cartInformation.UsedPoints) - discountAmount
	//不能减为负数，最低为0
	if totalAmount < 0 {
		totalAmount = 0
	}
	err = global.GVA_DB.Table("checkout_sessions").
		Where("user_id=?", userId).
		Update("total_amount", totalAmount).Error
	if err != nil {
		return nil, fmt.Errorf("更新优惠后的价格失败")
	}
	//更新checkout_sessions这个表里面的applied_coupon_id和coupon_discount_amount
	//先找到couponCode对应的couponId
	var coupinId string
	err = global.GVA_DB.Table("coupons").
		Select("id").
		Where("coupon_code=?", couponCode).
		Scan(&coupinId).Error
	if err != nil {
		return nil, err
	}

	err = global.GVA_DB.Table("checkout_sessions").
		Where("user_id=?", userId).
		Updates(map[string]interface{}{
			"applied_coupon_id":      coupinId,
			"coupon_discount_amount": discountAmount,
		}).Error
	if err != nil {
		return nil, err
	}

	return &dto.CheckoutInfoResponseV1_1{
		Message: "优惠券使用成功",
	}, nil
}

// 计算优惠金额
func (p *Payment) calculateDiscount(discountType string, discountValue, maxDiscountAmount, cartTotalAmount float64) float64 {
	var discountAmount float64

	if discountType == "percentage" {
		// 百分比优惠
		discountAmount = (discountValue / 100.0) * cartTotalAmount
	} else if discountType == "fixed" {
		// 固定金额优惠
		discountAmount = discountValue
	}

	// 限制最大优惠金额
	if discountAmount > maxDiscountAmount {
		discountAmount = maxDiscountAmount
	}
	if discountAmount > cartTotalAmount {
		discountAmount = cartTotalAmount
	}

	return discountAmount
}

// 解除coupon
func (p *Payment) RemoveCouponInfo(userId string, couponCode string) (*dto.CheckoutInfoResponseV1_1, error) {
	//同样，先查有没有这个优惠券
	var couponFind dto.AvailableCouponInfo
	err := global.GVA_DB.Table("coupons").
		Select(`
				id AS coupon_id,
				coupon_code,
				name,
				discount_value,
				discount_type,
				min_purchase_amount,
				max_discount_amount
				`).
		Where("coupon_code = ? AND NOW() BETWEEN start_date AND end_date", couponCode).
		Scan(&couponFind).Error
	if err != nil {
		return nil, fmt.Errorf("优惠券不存在或已过期")
	}
	//再查用户有没有使用这个coupon
	var usedCoupon int64
	err = global.GVA_DB.Table("user_coupons").
		Where("user_id=? AND coupon_id=? AND is_used=1", userId, couponFind.CouponID).
		Count(&usedCoupon).Error
	if err != nil {
		return nil, fmt.Errorf("查询状态失败")
	}
	if usedCoupon == 0 {
		return nil, fmt.Errorf("未使用的优惠券,无法解除")
	}
	//把is_userd变为0
	err = global.GVA_DB.Table("user_coupons").
		Where("user_id=? AND coupon_id=?", userId, couponFind.CouponID).
		Updates(map[string]interface{}{
			"is_used": 0,
			"used_at": nil,
		}).Error
	if err != nil {
		return nil, fmt.Errorf("撤销优惠券失败")
	}
	//计算取消之后的价格
	//查优惠前购物车总金额和被使用的积分
	var cartInformation struct {
		CartTotalAmount float64 `gorm:"column:cart_subtotal"`
		UsedPoints      int     `gorm:"column:used_points"`
		ShippingFee     float64 `gorm:"column:shipping_fee"`
	}
	err = global.GVA_DB.Table("checkout_sessions").
		Select("cart_subtotal,used_points,shipping_fee").
		Where("user_id=?", userId).
		Scan(&cartInformation).Error
	if err != nil {
		return nil, fmt.Errorf("获取购物车信息失败")
	}
	noDiscountAmount := cartInformation.CartTotalAmount - float64(cartInformation.UsedPoints) + cartInformation.ShippingFee
	//更新checkout_sessions表里关于coupon的状态
	err = global.GVA_DB.Table("checkout_sessions").
		Where("user_id=?", userId).
		Updates(map[string]interface{}{
			"applied_coupon_id":      nil,
			"coupon_discount_amount": 0,
			"total_amount":           noDiscountAmount,
		}).Error
	if err != nil {
		return nil, fmt.Errorf("更新金额失败")
	}
	return &dto.CheckoutInfoResponseV1_1{
		Message: "优惠券撤销成功",
	}, nil
}

func (p *Payment) UsePointInfo(userId string, point int) (*dto.CheckoutInfoResponseV1_1, error) {
	//先查用户积分
	var userPoint int
	err := global.GVA_DB.Table("user_points").
		Select("available_points").
		Where("user_id=?", userId).
		Scan(&userPoint).Error
	if err != nil {
		return nil, fmt.Errorf("查询积分失败")
	}
	if userPoint < point {
		return nil, fmt.Errorf("积分不足,无法使用")
	}

	//还是要先查checkout_sessions表下的数据
	var cartInformation struct {
		CartTotalAmount      float64 `gorm:"column:cart_subtotal"`
		ShippingFee          float64 `gorm:"column:shipping_fee"`
		CouponDiscountAmount float64 `gorm:"column:coupon_discount_amount"`
	}
	err = global.GVA_DB.Table("checkout_sessions").
		Select("cart_subtotal,coupon_discount_amount,shipping_fee").
		Where("user_id=?", userId).
		Scan(&cartInformation).Error
	if err != nil {
		return nil, fmt.Errorf("获取购物车信息失败")
	}
	pointAmount := cartInformation.CartTotalAmount - cartInformation.CouponDiscountAmount + cartInformation.ShippingFee - float64(point)
	//同样不能为负数
	if pointAmount < 0 {
		pointAmount = 0
	}
	//用积分就要更新checkout_sessions数据
	err = global.GVA_DB.Table("checkout_sessions").
		Where("user_id=?", userId).
		Updates(map[string]interface{}{
			"used_points":            point,
			"points_discount_amount": point,
			"total_amount":           pointAmount,
		}).Error
	if err != nil {
		return nil, fmt.Errorf("更新数据失败")
	}
	//减少用户的积分
	err = global.GVA_DB.Table("user_points").
		Where("user_id=?", userId).
		Update("available_points", userPoint-point).Error
	if err != nil {
		return nil, fmt.Errorf("扣除积分失败")
	}
	return &dto.CheckoutInfoResponseV1_1{
		Message: "积分使用成功",
	}, nil
}

// 撤销积分
func (p *Payment) RemovePointInfo(userId string) (*dto.CheckoutInfoResponseV1_1, error) {
	//先查是否用了积分和数据
	var cartInformation struct {
		CartSubtotal         float64 `json:"column:cart_subtotal"`
		ShippingFee          float64 `json:"column:shipping_fee"`
		CouponDiscountAmount float64 `json:"column:coupon_discount_amount"`
		PointsDiscountAmount float64 `json:"column:points_discount_amount"`
	}
	err := global.GVA_DB.Table("checkout_sessions").
		Select("cart_subtotal,coupon_discount_amount,points_discount_amount,shipping_fee").
		Where("user_id=?", userId).
		Scan(&cartInformation).Error
	if err != nil {
		return nil, fmt.Errorf("查询购物车数据失败")
	}
	if cartInformation.PointsDiscountAmount == 0 {
		return nil, fmt.Errorf("没有可以撤销的积分")
	}
	noPointsAmount := cartInformation.CartSubtotal - cartInformation.CouponDiscountAmount + cartInformation.ShippingFee
	if noPointsAmount < 0 {
		noPointsAmount = 0
	}
	//更新checkout_sessions下的数据
	err = global.GVA_DB.Table("checkout_sessions").
		Where("user_id=?", userId).
		Updates(map[string]interface{}{
			"used_points":            0,
			"points_discount_amount": 0.00,
			"total_amount":           noPointsAmount,
		}).Error
	if err != nil {
		return nil, fmt.Errorf("更新数据失败")
	}
	//把扣掉的积分再加回到账户上面
	err = global.GVA_DB.Table("user_points").
		Where("user_id=?", userId).
		Update("available_points", gorm.Expr("available_points + ?", int(cartInformation.PointsDiscountAmount))).Error
	if err != nil {
		return nil, fmt.Errorf("返还积分失败")
	}
	return &dto.CheckoutInfoResponseV1_1{
		Message: "积分撤销成功",
	}, nil
}
