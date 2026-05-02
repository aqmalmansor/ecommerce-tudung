package enums

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "ACTIVE"
	ProductStatusDraft    ProductStatus = "DRAFT"
	ProductStatusArchived ProductStatus = "ARCHIVED"
)

func (s ProductStatus) String() string {
	return string(s)
}

func (s ProductStatus) IsValid() bool {
	switch s {
	case ProductStatusActive, ProductStatusDraft, ProductStatusArchived:
		return true
	}
	return false
}
