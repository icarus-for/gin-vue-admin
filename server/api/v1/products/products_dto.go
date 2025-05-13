package products

import "time"

type ProductInfoResponse struct {
	ID                 string               `json:"id"`
	ProductCode        string               `json:"product_code,omitempty"`
	Name               string               `json:"name"`
	Description        *string              `json:"description,omitempty"`
	Status             string               `json:"status"`
	IsTaxable          bool                 `json:"is_taxable"`
	MetaTitle          *string              `json:"meta_title,omitempty"`
	MetaDescription    *string              `json:"meta_description,omitempty"`
	CreatedAtFormatted string               `json:"created_at_formatted"`
	UpdatedAtFormatted string               `json:"updated_at_formatted"`
	TargetSKUInfo      *TargetSKUInfo       `json:"target_sku_info,omitempty"` // 対象SKU情報 (Nullable)
	VariantOptions     []VariantOptionGroup `json:"variant_options,omitempty"` // 全バリエーション選択肢
	Category           *CategoryInfo        `json:"category,omitempty"`
	Brand              *BrandInfo           `json:"brand,omitempty"`
}

// TargetSKUInfo 対象SKUの詳細情報
type TargetSKUInfo struct {
	SkuID        string          `json:"sku_id"`
	SkuCode      string          `json:"sku_code,omitempty"`
	Price        *PriceInfo      `json:"price"`
	StockInfo    *StockInfo      `json:"stock_info"`
	PrimaryImage *ImageInfo      `json:"primary_image,omitempty"`
	Attributes   []AttributeInfo `json:"attributes"`
}

// PriceInfo 価格情報 (計算用数値と表示用文字列を含む)
type PriceInfo struct {
	Amount                  float64  `json:"amount"`
	FormattedAmount         string   `json:"formatted_amount"`
	Currency                string   `json:"currency"`
	Type                    string   `json:"type"`
	TypeName                string   `json:"type_name"`
	OriginalAmount          *float64 `json:"original_amount,omitempty"`
	FormattedOriginalAmount *string  `json:"formatted_original_amount,omitempty"`
}

// StockInfo 在庫状況概要
type StockInfo struct {
	Status           string  `json:"status"`
	StatusText       string  `json:"status_text"`
	DeliveryEstimate *string `json:"delivery_estimate,omitempty"`
}

// ImageInfo 画像情報
type ImageInfo struct {
	ID      int     `json:"id"`
	URL     string  `json:"url"`
	AltText *string `json:"alt_text,omitempty"`
}

// AttributeInfo 対象SKUの属性情報
type AttributeInfo struct {
	AttributeID   int     `json:"attribute_id"`
	AttributeName string  `json:"attribute_name"`
	OptionID      *int    `json:"option_id,omitempty"`
	OptionValue   *string `json:"option_value,omitempty"`
	ValueString   *string `json:"value_string,omitempty"`
	// ... 他の value_xxx 型
}

// VariantOptionGroup SKUバリエーション軸ごとの選択肢グループ
type VariantOptionGroup struct {
	AttributeID   int             `json:"attribute_id"`
	AttributeName string          `json:"attribute_name"`
	AttributeCode string          `json:"attribute_code"`
	DisplayType   string          `json:"display_type"` // 'image', 'text', etc.
	Options       []VariantOption `json:"options"`
}

// VariantOption 個々のバリエーション選択肢
type VariantOption struct {
	OptionID     int      `json:"option_id"`
	OptionValue  string   `json:"option_value"`
	OptionCode   string   `json:"option_code"`
	ImageURL     *string  `json:"image_url,omitempty"`      // Swatch画像等
	IsSelectable bool     `json:"is_selectable"`            // 在庫等による選択可否
	LinkedSkuIDs []string `json:"linked_sku_ids,omitempty"` // 関連SKU (任意)
}

// CategoryInfo カテゴリ情報
type CategoryInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Level    int    `json:"level"`
	ParentID *int   `json:"parent_id,omitempty"`
	// Breadcrumbs []Breadcrumb `json:"breadcrumbs,omitempty"`
}

// BrandInfo ブランド情報
type BrandInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse エラーレスポンス構造
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}

type Category struct {
	ID          int       `gorm:"primaryKey;comment:カテゴリID (手動割当)"`
	Name        string    `gorm:"size:255;not null;comment:カテゴリ名"`
	Description *string   `gorm:"type:text;comment:カテゴリ説明"`
	ParentID    *int      `gorm:"comment:親カテゴリID (NULLの場合はトップレベル)"`
	Level       int       `gorm:"not null;comment:階層レベル (0始まり)"`
	SortOrder   int       `gorm:"default:0;comment:表示順"`
	IsActive    bool      `gorm:"not null;default:true;comment:有効なカテゴリか"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:作成日時"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新日時"`

	Parent     *Category   `gorm:"foreignKey:ParentID"`
	Children   []Category  `gorm:"foreignKey:ParentID"`
	Products   []Product   `gorm:"foreignKey:CategoryID"`
	Attributes []Attribute `gorm:"many2many:category_attributes;foreignKey:ID;joinForeignKey:CategoryID;References:ID;joinReferences:AttributeID"`
}

// Attribute 属性定義
type Attribute struct {
	ID            int       `gorm:"primaryKey;comment:属性ID (手動割当)"`
	Name          string    `gorm:"size:255;not null;comment:属性名 (例: カラー, サイズ)"`
	AttributeCode string    `gorm:"size:100;not null;unique;comment:属性コード (例: color, size, material)"`
	InputType     string    `gorm:"type:enum('select','text','number','boolean','textarea');not null;comment:入力形式"`
	IsFilterable  bool      `gorm:"not null;default:false;comment:絞り込み検索対象か"`
	IsComparable  bool      `gorm:"not null;default:false;comment:商品比較対象か"`
	SortOrder     int       `gorm:"default:0;comment:表示順"`
	CreatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:作成日時"`
	UpdatedAt     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新日時"`

	Options    []AttributeOption `gorm:"foreignKey:AttributeID"`
	Categories []Category        `gorm:"many2many:category_attributes;foreignKey:ID;joinForeignKey:AttributeID;References:ID;joinReferences:CategoryID"`
	SkuValues  []SkuValue        `gorm:"foreignKey:AttributeID"`
}

// AttributeOption 属性選択肢
type AttributeOption struct {
	ID          int       `gorm:"primaryKey;comment:属性選択肢ID (手動割当)"`
	AttributeID int       `gorm:"not null;comment:属性ID (input_typeがselectの場合)"`
	Value       string    `gorm:"size:255;not null;comment:表示値 (例: レッド, Mサイズ)"`
	OptionCode  string    `gorm:"size:100;not null;comment:選択肢コード (例: red, size_m)"`
	SortOrder   int       `gorm:"default:0;comment:表示順"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:作成日時"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新日時"`

	Attribute Attribute  `gorm:"foreignKey:AttributeID"`
	SkuValues []SkuValue `gorm:"foreignKey:OptionID"`
}

// PriceType 価格種別マスタ
type PriceType struct {
	ID       int    `gorm:"primaryKey;comment:価格種別ID (手動割当)"`
	TypeCode string `gorm:"size:50;not null;unique;comment:価格種別コード (例: regular, sale, member_special)"`
	Name     string `gorm:"size:100;not null;comment:価格種別名"`

	Prices []Price `gorm:"foreignKey:PriceTypeID"`
}

// InventoryLocation 在庫拠点マスタ
type InventoryLocation struct {
	ID           int    `gorm:"primaryKey;comment:在庫拠点ID (手動割当)"`
	LocationCode string `gorm:"size:50;not null;unique;comment:拠点コード (例: WAREHOUSE_EAST, STORE_SHIBUYA)"`
	Name         string `gorm:"size:100;not null;comment:拠点名"`
	LocationType string `gorm:"type:enum('warehouse','store','distribution_center');not null;comment:拠点タイプ"`

	Inventories []Inventory `gorm:"foreignKey:LocationID"`
}

// SalesChannel 販売チャネルマスタ
type SalesChannel struct {
	ID          int    `gorm:"primaryKey;comment:販売チャネルID (手動割当)"`
	ChannelCode string `gorm:"size:50;not null;unique;comment:チャネルコード (例: ONLINE_JP, STORE)"`
	Name        string `gorm:"size:100;not null;comment:チャネル名"`

	SkuAvailabilities []SkuAvailability `gorm:"foreignKey:SalesChannelID"`
}

// Product 商品の基本情報 (SKUの親)
type Product struct {
	ID              string     `gorm:"primaryKey;type:char(36);comment:商品ID (UUID)"`
	Name            string     `gorm:"size:255;not null;comment:商品名"`
	Description     *string    `gorm:"type:text;comment:商品説明"`
	ProductCode     *string    `gorm:"size:100;unique;comment:商品管理コード (ニトリの品番など)"`
	CategoryID      int        `gorm:"not null;comment:主カテゴリID"`
	BrandID         *int       `gorm:"comment:ブランドID (将来用)"`
	DefaultSkuID    *string    `gorm:"type:char(36);comment:代表SKU ID (初期表示用)"`
	Status          string     `gorm:"type:enum('draft','active','inactive','discontinued');not null;default:'draft';comment:商品ステータス"`
	IsTaxable       bool       `gorm:"not null;default:true;comment:課税対象か"`
	MetaTitle       *string    `gorm:"size:255;comment:SEO用タイトル"`
	MetaDescription *string    `gorm:"size:500;comment:SEO用説明文"`
	CreatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP;comment:作成日時"`
	UpdatedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新日時"`
	DeletedAt       *time.Time `gorm:"index;comment:論理削除日時"`

	Category   Category     `gorm:"foreignKey:CategoryID"`
	DefaultSku *ProductSku  `gorm:"foreignKey:DefaultSkuID"`
	Skus       []ProductSku `gorm:"foreignKey:ProductID"`
}

// ProductSku SKU (在庫管理の最小単位)
type ProductSku struct {
	ID        string     `gorm:"primaryKey;type:char(36);comment:SKU ID (UUID)"`
	ProductID string     `gorm:"type:char(36);not null;comment:商品ID"`
	SkuCode   *string    `gorm:"size:150;unique;comment:SKUコード (商品コード + バリエーション識別子)"`
	Status    string     `gorm:"type:enum('active','inactive','discontinued');not null;default:'active';comment:SKUステータス"`
	Barcode   *string    `gorm:"size:50;comment:バーコード (JAN/UPCなど)"`
	Weight    *float64   `gorm:"type:decimal(10,3);comment:重量 (kg)"`
	Width     *float64   `gorm:"type:decimal(10,2);comment:幅 (cm)"`
	Height    *float64   `gorm:"type:decimal(10,2);comment:高さ (cm)"`
	Depth     *float64   `gorm:"type:decimal(10,2);comment:奥行 (cm)"`
	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP;comment:作成日時"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新日時"`
	DeletedAt *time.Time `gorm:"index;comment:論理削除日時"`

	Product        Product           `gorm:"foreignKey:ProductID"`
	Values         []SkuValue        `gorm:"foreignKey:SkuID"`
	Images         []SkuImage        `gorm:"foreignKey:SkuID"`
	Prices         []Price           `gorm:"foreignKey:SkuID"`
	Inventories    []Inventory       `gorm:"foreignKey:SkuID"`
	Availabilities []SkuAvailability `gorm:"foreignKey:SkuID"`
}

// CategoryAttribute カテゴリ別利用属性紐付け
type CategoryAttribute struct {
	CategoryID         int  `gorm:"primaryKey;comment:カテゴリID"`
	AttributeID        int  `gorm:"primaryKey;comment:属性ID"`
	IsRequired         bool `gorm:"not null;default:false;comment:SKU定義に必須か"`
	IsVariantAttribute bool `gorm:"not null;default:false;comment:SKUのバリエーション軸となる属性か (例: 色、サイズ)"`
	SortOrder          int  `gorm:"default:0;comment:カテゴリ内での属性表示順"`

	Category  Category  `gorm:"foreignKey:CategoryID"`
	Attribute Attribute `gorm:"foreignKey:AttributeID"`
}

// SkuValue SKU属性値
type SkuValue struct {
	ID           int64    `gorm:"primaryKey;comment:SKU属性値ID (手動割当)"`
	SkuID        string   `gorm:"type:char(36);not null;comment:SKU ID"`
	AttributeID  int      `gorm:"not null;comment:属性ID"`
	OptionID     *int     `gorm:"comment:選択肢ID (input_type=select)"`
	ValueString  *string  `gorm:"size:255;comment:文字列値 (input_type=text)"`
	ValueNumber  *float64 `gorm:"type:decimal(15,4);comment:数値 (input_type=number)"`
	ValueBoolean *bool    `gorm:"comment:真偽値 (input_type=boolean)"`
	ValueText    *string  `gorm:"type:text;comment:長文テキスト値 (input_type=textarea)"`

	Sku       ProductSku       `gorm:"foreignKey:SkuID"`
	Attribute Attribute        `gorm:"foreignKey:AttributeID"`
	Option    *AttributeOption `gorm:"foreignKey:OptionID"`
}

// SkuImage SKU画像
type SkuImage struct {
	ID        int       `gorm:"primaryKey;comment:画像ID (手動割当)"`
	SkuID     string    `gorm:"type:char(36);not null;comment:SKU ID"`
	ImageURL  string    `gorm:"size:500;not null;comment:画像URL (CDNなど)"`
	AltText   *string   `gorm:"size:255;comment:代替テキスト"`
	SortOrder int       `gorm:"default:0;comment:表示順"`
	ImageType string    `gorm:"type:enum('main','swatch','gallery','detail');default:'gallery';comment:画像タイプ"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;comment:作成日時"`

	Sku ProductSku `gorm:"foreignKey:SkuID"`
}

// Price SKU価格
type Price struct {
	ID           int64      `gorm:"primaryKey;comment:価格ID (手動割当)"`
	SkuID        string     `gorm:"type:char(36);not null;comment:SKU ID"`
	PriceTypeID  int        `gorm:"not null;comment:価格種別ID"`
	Price        float64    `gorm:"type:decimal(12,2);not null;comment:価格"`
	CurrencyCode string     `gorm:"size:3;not null;default:'JPY';comment:通貨コード"`
	StartDate    *time.Time `gorm:"comment:適用開始日時"`
	EndDate      *time.Time `gorm:"comment:適用終了日時"`
	IsActive     bool       `gorm:"not null;default:true;comment:有効な価格設定か"`
	CreatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP;comment:作成日時"`
	UpdatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新日時"`

	Sku       ProductSku `gorm:"foreignKey:SkuID"`
	PriceType PriceType  `gorm:"foreignKey:PriceTypeID"`
}

// Inventory 在庫
type Inventory struct {
	ID               int64     `gorm:"primaryKey;comment:在庫ID (手動割当)"`
	SkuID            string    `gorm:"type:char(36);not null;comment:SKU ID"`
	LocationID       int       `gorm:"not null;comment:在庫拠点ID"`
	Quantity         int       `gorm:"not null;default:0;comment:物理在庫数"`
	ReservedQuantity int       `gorm:"not null;default:0;comment:引当済在庫数"`
	LastUpdated      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:最終更新日時"`

	Sku      ProductSku        `gorm:"foreignKey:SkuID"`
	Location InventoryLocation `gorm:"foreignKey:LocationID"`
}

func (Inventory) TableName() string {
	return "inventory"
}

// SkuAvailability SKU販売可否
type SkuAvailability struct {
	ID             int64      `gorm:"primaryKey;comment:販売可否ID (手動割当)"`
	SkuID          string     `gorm:"type:char(36);not null;comment:SKU ID"`
	SalesChannelID int        `gorm:"not null;comment:販売チャネルID"`
	IsAvailable    bool       `gorm:"not null;default:true;comment:販売可能か"`
	AvailableFrom  *time.Time `gorm:"comment:販売開始日時"`
	AvailableUntil *time.Time `gorm:"comment:販売終了日時"`

	Sku          ProductSku   `gorm:"foreignKey:SkuID"`
	SalesChannel SalesChannel `gorm:"foreignKey:SalesChannelID"`
}

func (SkuAvailability) TableName() string {
	return "sku_availability"
}
