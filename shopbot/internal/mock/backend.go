package mock

import (
	"context"
	"fmt"
	"sync"

	"shopbot/internal/backend"
)

type Backend struct {
	seed  Seed
	mu    sync.Mutex
	carts map[int64][]cartLine
}

type cartLine struct {
	productID string
	Quantity  int
}

var _ backend.Backend = (*Backend)(nil)

func New(seed Seed) *Backend {
	return &Backend{
		seed:  seed,
		carts: make(map[int64][]cartLine),
	}
}

func (b *Backend) Cart(ctx context.Context, userID int64) (backend.Cart, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var cart backend.Cart

	for _, line := range b.carts[userID] {
		p, err := b.Product(ctx, line.productID)
		if err != nil {
			return backend.Cart{}, err
		}

		cart.Items = append(cart.Items, backend.CartItem{
			Product: p,
			Qty:     line.Quantity,
		})
		cart.Total += p.Price * int64(line.Quantity)
	}

	return cart, nil
}

func (b *Backend) AddToCart(ctx context.Context, userID int64, productID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, err := b.Product(ctx, productID)
	if err != nil {
		return err
	}

	for i := range b.carts[userID] {
		if b.carts[userID][i].productID == productID {
			b.carts[userID][i].Quantity++
			return nil
		}
	}

	b.carts[userID] = append(b.carts[userID], cartLine{
		productID: productID,
		Quantity:  1,
	})

	return nil
}

func (b *Backend) ChangeQty(ctx context.Context, userID int64, productID string, delta int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	lines := b.carts[userID]

	for i, line := range lines {
		if line.productID == productID {
			lines[i].Quantity += delta

			if lines[i].Quantity <= 0 {
				newLines := []cartLine{}

				for _, item := range lines {
					if item.productID != productID {
						newLines = append(newLines, item)
					}
				}

				b.carts[userID] = newLines
			}

			return nil
		}

	}

	return backend.ErrNotFound
}

func (b *Backend) RemoveFromCart(ctx context.Context, userID int64, productID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	lines := b.carts[userID]
	newLines := []cartLine{}
	found := false

	for _, line := range lines {
		if line.productID == productID {
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}

	if !found {
		return backend.ErrNotFound
	}

	b.carts[userID] = newLines
	return nil
}

func (b *Backend) ClearCart(ctx context.Context, userID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.carts, userID)

	return nil
}

func (b *Backend) Categories(ctx context.Context) ([]backend.Category, error) {
	return b.seed.Categories, nil
}

func (b *Backend) ProductsByCategory(ctx context.Context, categoryID string) ([]backend.Product, error) {
	var result []backend.Product
	for _, p := range b.seed.Products {
		if p.CategoryID == categoryID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (b *Backend) Product(ctx context.Context, id string) (backend.Product, error) {
	for _, p := range b.seed.Products {
		if p.ID == id {
			return p, nil
		}
	}
	return backend.Product{}, fmt.Errorf("product %q: %w", id, backend.ErrNotFound)
}

func (b *Backend) PickupAddress(ctx context.Context) (string, error) {
	return b.seed.PickupAddress, nil
}
