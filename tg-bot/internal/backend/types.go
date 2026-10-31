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

// DeliveryMethod is how the customer gets the order.
type DeliveryMethod string

const (
	Pickup   DeliveryMethod = "pickup"   // the customer picks it up at the school
	Delivery DeliveryMethod = "delivery" // we bring it to the customer's location
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// CheckoutInfo is what the user enters at checkout; CreateOrder will receive it.
type CheckoutInfo struct {
	Method   DeliveryMethod `json:"method"`
	Location *Location      `json:"location,omitempty"` // nil for Pickup
	Phone    string         `json:"phone"`
	Time     string         `json:"time"` // desired time, as the user wrote it
}
