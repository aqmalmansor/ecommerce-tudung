package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint             `gorm:"primarykey;autoIncrement" json:"id"`
	Title       string           `gorm:"not null" json:"title"`
	Description string           `json:"description"`
	Handle      string           `gorm:"uniqueIndex;not null" json:"handle"`
	Vendor      string           `json:"vendor"`
	ProductType string           `json:"product_type"`
	Status      string           `gorm:"default:draft" json:"status"`
	Tags        string           `json:"tags"`
	Images      []ProductImage   `gorm:"foreignKey:ProductID" json:"images,omitempty"`
	Options     []ProductOption  `gorm:"foreignKey:ProductID" json:"options,omitempty"`
	Variants    []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
}

type ProductOption struct {
	ID        uint                 `gorm:"primarykey;autoIncrement" json:"id"`
	ProductID uint                 `gorm:"not null;index" json:"product_id"`
	Name      string               `gorm:"not null" json:"name"`
	Position  int                  `gorm:"default:1" json:"position"`
	Values    []ProductOptionValue `gorm:"foreignKey:OptionID" json:"values,omitempty"`
}

type ProductOptionValue struct {
	ID       uint   `gorm:"primarykey;autoIncrement" json:"id"`
	OptionID uint   `gorm:"not null;index" json:"option_id"`
	Value    string `gorm:"not null" json:"value"`
}

type ProductVariant struct {
	ID               uint           `gorm:"primarykey;autoIncrement" json:"id"`
	ProductID        uint           `gorm:"not null;index" json:"product_id"`
	Title            string         `json:"title"`
	SKU              string         `gorm:"uniqueIndex" json:"sku"`
	Barcode          string         `json:"barcode"`
	Price            float64        `gorm:"not null" json:"price"`
	CompareAtPrice   *float64       `json:"compare_at_price"`
	InventoryQty     int            `gorm:"default:0" json:"inventory_quantity"`
	Option1          string         `json:"option1"`
	Option2          string         `json:"option2"`
	Option3          string         `json:"option3"`
	Weight           float64        `json:"weight"`
	RequiresShipping bool           `gorm:"default:true" json:"requires_shipping"`
	Taxable          bool           `gorm:"default:true" json:"taxable"`
	ImageID          *uint          `json:"image_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProductImage struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	Src       string    `gorm:"not null" json:"src"`
	AltText   string    `json:"alt_text"`
	Position  int       `gorm:"default:1" json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
