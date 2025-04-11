package domain

type Product struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	StockLevel int     `json:"stock_level"`
	Price      float64 `json:"price"`
}

type ProductRepository interface {
	Create(product Product) error
	GetByID(id string) (Product, error)
	Update(id string, product Product) error
	Delete(id string) error
	List() ([]Product, error)
}
