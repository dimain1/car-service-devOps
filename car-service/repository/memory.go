package repository

import (
	"context"
	"sync"

	"dmitrii.com/car-service/car"
)

type memoryRepo struct {
	mu   sync.RWMutex
	cars map[string]car.Car
}

// NewMemoryRepository создаёт in-memory хранилище заказов
func NewMemoryRepository() car.Repository {
	return &memoryRepo{
		cars: make(map[string]car.Car),
	}
}
func (r *memoryRepo) CreateCar(ctx context.Context, c car.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cars[c.ID] = c
	return nil
}
func (r *memoryRepo) GetCarByID(ctx context.Context, id string) (car.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.cars[id]
	if !ok {
		return car.Car{}, car.ErrCarNotFound
	}
	return c, nil
}
func (r *memoryRepo) MarkCarSold(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.cars[id]
	if !ok {
		return car.ErrCarNotFound
	}
	c.Sold = true
	r.cars[id] = c
	return nil
}
