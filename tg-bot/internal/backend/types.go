package backend

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	PhotoURL    string `json:"photo_url"`
}

type CartItem struct {
	Product Product
	Qty     int
}

type Cart struct {
	Items []CartItem
	Total int64
}
