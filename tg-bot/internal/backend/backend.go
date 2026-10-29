package backend

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Backend interface {
	Categories(ctx context.Context) ([]Category, error)
	ProductsByCategory(ctx context.Context, categoryID string) ([]Product, error)
	Product(ctx context.Context, id string) (Product, error)
	PickupAddress(ctx context.Context) (string, error)

	Cart(ctx context.Context, userID int64) (Cart, error)
	AddToCart(ctx context.Context, userID int64, ProductID string) error
	ChangeQty(ctx context.Context, userID int64, ProductID string, delta int) error
	RemoveFromCart(ctx context.Context, userID int64, ProductID string) error
	ClearCart(ctx context.Context, userID int64) error
}
