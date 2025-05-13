package dto

type PaymentMethodInfo struct {
	MethodID    int     `json:"method_id"`
	MethodCode  string  `json:"method_code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// CheckoutInfoResponseV1_1 クーポン・ポイント情報取得APIのルートレスポンス (修正版)
type CheckoutInfoResponseV1_1 struct {
	AvailableCoupons     []AvailableCouponInfo `json:"available_coupons,omitempty"`      // 利用可能なクーポンリスト
	UserPoints           *UserPointInfo        `json:"user_points,omitempty"`            // 保有ポイント情報
	CurrentCheckoutState *CurrentCheckoutState `json:"current_checkout_state,omitempty"` // ★現在のチェックアウト状態
	Message              string                `json:"message,omitempty"`                // メッセージ
}

// AvailableCouponInfo 利用可能なクーポン情報
type AvailableCouponInfo struct {
	CouponID          uint64  `json:"coupon_id"`
	CouponCode        string  `json:"coupon_code"`
	Name              string  `json:"name"`
	Description       *string `json:"description,omitempty"`
	DiscountText      string  `json:"discount_text"` // 例: "10% OFF (最大2,000円引)", "500円引き"
	DiscountValue     float64 `json:"discount_value,omitempty"`
	DiscountType      string  `json:"discount_type,omitempty"`
	MinPurchaseAmount float64 `json:"min_purchase_amount,omitempty"`
	MaxDiscountAmount float64 `json:"max_discount_amount,omitempty"`
	ShippingFee       float64 `json:"shipping_fee,omitempty"`
}

// UserPointInfo ユーザー保有ポイント情報
type UserPointInfo struct {
	AvailablePoints int `json:"available_points"`
}

// CurrentCheckoutState 現在のチェックアウト状態を表すDTO
type CurrentCheckoutState struct {
	//CouponID                      uint64             `json:"coupon_id"`
	CartSubtotalAmountFormatted   string             `json:"cart_subtotal_formatted"`          // カート商品小計 (割引前、表示用)
	AppliedCouponInfo             *AppliedCouponInfo `json:"applied_coupon_info,omitempty"`    // 適用中クーポン情報 (Nullable)
	CouponDiscountAmountFormatted string             `json:"coupon_discount_amount_formatted"` // クーポン割引額 (表示用)
	UsedPoints                    int                `json:"used_points"`                      // 利用ポイント数
	PointsDiscountAmountFormatted string             `json:"points_discount_amount_formatted"` // ポイント割引額 (表示用)
	ShippingFeeFormatted          string             `json:"shipping_fee_formatted"`           // 送料 (表示用、別途計算の場合あり)
	TotalAmountFormatted          string             `json:"total_amount_formatted"`           // ★最終支払総額 (表示用)
	//DiscountAmount                float64            `json:"discount_amount"`

	// 内部計算用の数値も保持 (JSONには含めないか、開発用に含めるかは選択)
	CartSubtotalAmount   float64 `json:"-"`
	CouponDiscountAmount float64 `json:"-"`
	PointsDiscountAmount float64 `json:"-"`
	ShippingFee          float64 `json:"-"`
	TotalAmount          float64 `json:"-"`
}

// AppliedCouponInfo 適用中クーポン情報
type AppliedCouponInfo struct {
	CouponID                uint64  `json:"coupon_id"`
	CouponCode              string  `json:"coupon_code"`
	Name                    string  `json:"name"`
	DiscountAmount          float64 `json:"discount_amount"`           // この注文での実際の割引額 (計算用)
	FormattedDiscountAmount string  `json:"formatted_discount_amount"` // 表示用
}

// ApplyCouponRequest クーポン適用APIのリクエストボディ
type ApplyCouponRequest struct {
	CouponCode string `json:"coupon_code" binding:"required"`
}

// UsePointsRequest ポイント利用APIのリクエストボディ
type UsePointsRequest struct {
	PointsToUse int `json:"points_to_use" binding:"required,min=1"`
}
