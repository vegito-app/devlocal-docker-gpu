package v1

type Order struct {
	ID           string                   `json:"id,omitempty"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	SaleType     string                   `json:"saleType"`
	WeightGrams  int                      `json:"weightGrams"`
	PriceCents   int                      `json:"priceCents"`
	Images       []VegetableImageResponse `json:"images"`
	CreatedAt    int64                    `json:"createdAt"`
	Status       string                   `json:"status"`
	VegetableIds []string                 `json:"vegetableIds"`
}
