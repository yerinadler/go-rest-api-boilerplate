package events

type ProductCreated struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UnitPrice   int32  `json:"unitPrice"`
}
