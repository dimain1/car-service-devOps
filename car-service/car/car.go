package car

import "context"

// Order представляет заказ
type Car struct {
	ID    string  `json:"id,omitempty"`
	Brand string  `json:"brand"`
	Model string  `json:"model"`
	Year  int64   `json:"year"`
	Price float64 `json:"price"`
	Sold  bool    `json:"sold"`
}

// Service описывает бизнес-логику сервиса заказов
type Service interface {
	Create(ctx context.Context, order Car) (string, error)
	GetByID(ctx context.Context, id string) (Car, error)
	MarkSold(ctx context.Context, id string) error
}

// Repository описывает хранилище заказов
type Repository interface {
	CreateCar(ctx context.Context, car Car) error
	GetCarByID(ctx context.Context, id string) (Car, error)
	MarkCarSold(ctx context.Context, id string) error
}
