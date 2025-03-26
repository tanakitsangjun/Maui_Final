package dto

type ProductSearchRequest struct {
	Keyword  string  `json:"keyword"`
	MinPrice float64 `json:"min_price"`
	MaxPrice float64 `json:"max_price"`
}

type AddToCartRequest struct {
	CartName  string `json:"cart_name"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CartItemResponse struct {
	CartItemID  int     `json:"cart_item_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       string  `json:"price"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
}

type CartResponse struct {
	CartID   int                `json:"cart_id"`
	CartName string             `json:"cart_name"`
	Items    []CartItemResponse `json:"items"`
}

type AllCartsResponse struct {
	CustomerID   int            `json:"customer_id"`
	CustomerName string         `json:"customer_name"`
	Carts        []CartResponse `json:"carts"`
}
