package handlers

import (
	"net/http"
	"strings"

	"be/enums"
	"be/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

type CreateProductRequest struct {
	Title       string                 `json:"title" binding:"required"`
	Description string                 `json:"description"`
	Handle      string                 `json:"handle" binding:"required"`
	Vendor      string                 `json:"vendor"`
	ProductType string                 `json:"product_type"`
	Status      string                 `json:"status"`
	Tags        string                 `json:"tags"`
	Options     []CreateOptionRequest  `json:"options"`
	Variants    []CreateVariantRequest `json:"variants"`
	Images      []CreateImageRequest   `json:"images"`
}

type CreateOptionRequest struct {
	Name     string   `json:"name" binding:"required"`
	Position int      `json:"position"`
	Values   []string `json:"values"`
}

type CreateVariantRequest struct {
	Title            string   `json:"title"`
	SKU              string   `json:"sku" binding:"required"`
	Barcode          string   `json:"barcode"`
	Price            float64  `json:"price" binding:"required"`
	CompareAtPrice   *float64 `json:"compare_at_price"`
	InventoryQty     int      `json:"inventory_quantity"`
	Option1          string   `json:"option1"`
	Option2          string   `json:"option2"`
	Option3          string   `json:"option3"`
	Weight           float64  `json:"weight"`
	RequiresShipping *bool    `json:"requires_shipping"`
	Taxable          *bool    `json:"taxable"`
}

type CreateImageRequest struct {
	Src      string `json:"src" binding:"required"`
	AltText  string `json:"alt_text"`
	Position int    `json:"position"`
}

type UpdateVariantRequest struct {
	Title          string   `json:"title"`
	SKU            string   `json:"sku"`
	Barcode        string   `json:"barcode"`
	Price          *float64 `json:"price"`
	CompareAtPrice *float64 `json:"compare_at_price"`
	InventoryQty   *int     `json:"inventory_quantity"`
	Option1        string   `json:"option1"`
	Option2        string   `json:"option2"`
	Option3        string   `json:"option3"`
	Weight         *float64 `json:"weight"`
	ImageID        *uint    `json:"image_id"`
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	var products []models.Product
	query := h.DB.Preload("Images").Preload("Options.Values").Preload("Variants")

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", enums.ProductStatusActive.String())
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := h.DB.Preload("Images").Preload("Options.Values").Preload("Variants").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := enums.ProductStatusDraft
	if req.Status != "" {
		status = enums.ProductStatus(req.Status)
		if !status.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be active, draft, or archived"})
			return
		}
	}

	handle := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Handle), " ", "-"))

	product := models.Product{
		Title:       req.Title,
		Description: req.Description,
		Handle:      handle,
		Vendor:      req.Vendor,
		ProductType: req.ProductType,
		Status:      status.String(),
		Tags:        req.Tags,
	}

	if err := h.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product - " + err.Error()})
		return
	}

	for i, opt := range req.Options {
		pos := opt.Position
		if pos == 0 {
			pos = i + 1
		}
		option := models.ProductOption{ProductID: product.ID, Name: opt.Name, Position: pos}
		if err := h.DB.Create(&option).Error; err != nil {
			continue
		}
		for _, val := range opt.Values {
			h.DB.Create(&models.ProductOptionValue{OptionID: option.ID, Value: val})
		}
	}

	for _, v := range req.Variants {
		requiresShipping, taxable := true, true
		if v.RequiresShipping != nil {
			requiresShipping = *v.RequiresShipping
		}
		if v.Taxable != nil {
			taxable = *v.Taxable
		}
		h.DB.Create(&models.ProductVariant{
			ProductID:        product.ID,
			Title:            v.Title,
			SKU:              v.SKU,
			Barcode:          v.Barcode,
			Price:            v.Price,
			CompareAtPrice:   v.CompareAtPrice,
			InventoryQty:     v.InventoryQty,
			Option1:          v.Option1,
			Option2:          v.Option2,
			Option3:          v.Option3,
			Weight:           v.Weight,
			RequiresShipping: requiresShipping,
			Taxable:          taxable,
		})
	}

	for i, img := range req.Images {
		pos := img.Position
		if pos == 0 {
			pos = i + 1
		}
		h.DB.Create(&models.ProductImage{ProductID: product.ID, Src: img.Src, AltText: img.AltText, Position: pos})
	}

	h.DB.Preload("Images").Preload("Options.Values").Preload("Variants").First(&product, product.ID)
	c.JSON(http.StatusCreated, gin.H{"product": product})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := h.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if status, ok := updates["status"].(string); ok {
		if !enums.ProductStatus(status).IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
	}

	if err := h.DB.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	h.DB.Preload("Images").Preload("Options.Values").Preload("Variants").First(&product, product.ID)
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := h.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err := h.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (h *ProductHandler) AddVariant(c *gin.Context) {
	productID := c.Param("id")
	var product models.Product
	if err := h.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var req CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requiresShipping, taxable := true, true
	if req.RequiresShipping != nil {
		requiresShipping = *req.RequiresShipping
	}
	if req.Taxable != nil {
		taxable = *req.Taxable
	}

	variant := models.ProductVariant{
		ProductID:        product.ID,
		Title:            req.Title,
		SKU:              req.SKU,
		Barcode:          req.Barcode,
		Price:            req.Price,
		CompareAtPrice:   req.CompareAtPrice,
		InventoryQty:     req.InventoryQty,
		Option1:          req.Option1,
		Option2:          req.Option2,
		Option3:          req.Option3,
		Weight:           req.Weight,
		RequiresShipping: requiresShipping,
		Taxable:          taxable,
	}

	if err := h.DB.Create(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create variant - " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"variant": variant})
}

func (h *ProductHandler) UpdateVariant(c *gin.Context) {
	productID := c.Param("id")
	variantID := c.Param("variantId")

	var variant models.ProductVariant
	if err := h.DB.Where("id = ? AND product_id = ?", variantID, productID).First(&variant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	var req UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.SKU != "" {
		updates["sku"] = req.SKU
	}
	if req.Barcode != "" {
		updates["barcode"] = req.Barcode
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.CompareAtPrice != nil {
		updates["compare_at_price"] = *req.CompareAtPrice
	}
	if req.InventoryQty != nil {
		updates["inventory_qty"] = *req.InventoryQty
	}
	if req.Option1 != "" {
		updates["option1"] = req.Option1
	}
	if req.Option2 != "" {
		updates["option2"] = req.Option2
	}
	if req.Option3 != "" {
		updates["option3"] = req.Option3
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.ImageID != nil {
		updates["image_id"] = *req.ImageID
	}

	if err := h.DB.Model(&variant).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"variant": variant})
}

func (h *ProductHandler) DeleteVariant(c *gin.Context) {
	productID := c.Param("id")
	variantID := c.Param("variantId")

	var variant models.ProductVariant
	if err := h.DB.Where("id = ? AND product_id = ?", variantID, productID).First(&variant).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
		return
	}

	if err := h.DB.Delete(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Variant deleted"})
}
