package repository

import (
	"errors"

	"github.com/mikestefanello/otcscanner/models"
)

// ErrNotFound is an error that indicates an order could not be found
var ErrNotFound = errors.New("Order not found")

// OrderRepository provides an interface for order repositories
type OrderRepository interface {
	// LoadByID loads an order with a given ID
	LoadByID(id string) (*models.Order, error)

	// LoadAll loads all orders
	LoadAll() (*models.Orders, error)

	// LoadCompleted loads completed orders
	LoadCompleted() (*models.Orders, error)

	// LoadIncomplete loads incomplete orders
	LoadIncomplete() (*models.Orders, error)

	// DeleteAll deletes all orders
	DeleteAll() error

	// DeleteCompleted deletes completed orders
	DeleteCompleted() error

	// UpdateOne updates a given order
	UpdateOne(order *models.Order) error

	// InsertMany inserts multiple new orders
	InsertMany(orders *models.Orders) error

	// CountAll counts all orders
	CountAll() (int64, error)

	// CountCompleted counts completed orders
	CountCompleted() (int64, error)

	// CountIncomplete counts incomplete orders
	CountIncomplete() (int64, error)
}
