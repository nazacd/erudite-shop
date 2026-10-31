package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"tg-bot/internal/backend"
)

type Backend struct {
	seed      Seed
	statePath string

	mu    sync.Mutex
	carts map[int64][]cartLine
}

// cartLine fields are exported so encoding/json can save them.
type cartLine struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// state is what is stored in the state file.
type state struct {
	Carts map[int64][]cartLine `json:"carts"`
}

var _ backend.Backend = (*Backend)(nil)

// New creates the mock and loads the saved state from statePath.
// If the file does not exist yet, it starts with empty carts.
func New(seed Seed, statePath string) (*Backend, error) {
	b := &Backend{
		seed:      seed,
		statePath: statePath,
		carts:     make(map[int64][]cartLine),
	}

	data, err := os.ReadFile(statePath)
	if os.IsNotExist(err) {
		return b, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", statePath, err)
	}
	var st state
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("bad JSON in %s: %w", statePath, err)
	}
	if st.Carts != nil {
		b.carts = st.Carts
	}
	return b, nil
}

// save writes the state to the file. b.mu must already be locked by the caller.
func (b *Backend) save() error {
	data, err := json.MarshalIndent(state{Carts: b.carts}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}
	// Write to a temporary file first and then rename it over the real one,
	// so a crash in the middle of writing never leaves a half-written file.
	tmp := b.statePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write state: %w", err)
	}
	if err := os.Rename(tmp, b.statePath); err != nil {
		return fmt.Errorf("write state: %w", err)
	}
	return nil
}

func (b *Backend) Cart(ctx context.Context, userID int64) (backend.Cart, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	var cart backend.Cart

	for _, line := range b.carts[userID] {
		p, err := b.Product(ctx, line.ProductID)
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
		if b.carts[userID][i].ProductID == productID {
			b.carts[userID][i].Quantity++
			return b.save()
		}
	}

	b.carts[userID] = append(b.carts[userID], cartLine{
		ProductID: productID,
		Quantity:  1,
	})

	return b.save()
}

func (b *Backend) ChangeQty(ctx context.Context, userID int64, productID string, delta int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	lines := b.carts[userID]

	for i, line := range lines {
		if line.ProductID == productID {
			lines[i].Quantity += delta

			if lines[i].Quantity <= 0 {
				newLines := []cartLine{}

				for _, item := range lines {
					if item.ProductID != productID {
						newLines = append(newLines, item)
					}
				}

				b.carts[userID] = newLines
			}

			return b.save()
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
		if line.ProductID == productID {
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}

	if !found {
		return backend.ErrNotFound
	}

	b.carts[userID] = newLines
	return b.save()
}

func (b *Backend) ClearCart(ctx context.Context, userID int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.carts, userID)

	return b.save()
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

// DeliveryEstimate ignores the location and always returns a flat price.
// The real backend will ask Yandex Taxi.
func (b *Backend) DeliveryEstimate(ctx context.Context, loc backend.Location) (int64, error) {
	return 20000, nil
}
