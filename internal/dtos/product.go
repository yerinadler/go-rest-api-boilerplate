package dtos

type ProductDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UnitPrice   int32  `json:"unitPrice"`
}
