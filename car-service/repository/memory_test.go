package repository

import (
	"context"
	"testing"

	"dmitrii.com/car-service/car"
)

func TestMemoryRepository(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	testCar := car.Car{
		ID:    "test-id",
		Brand: "Toyota",
		Model: "Camry",
		Year:  2022,
		Price: 25000,
	}

	// Создание автомобиля.
	if err := repo.CreateCar(ctx, testCar); err != nil {
		t.Fatalf("CreateCar() error: %v", err)
	}

	// Получение созданного автомобиля.
	got, err := repo.GetCarByID(ctx, testCar.ID)
	if err != nil {
		t.Fatalf("GetCarByID() error: %v", err)
	}
	if got.Brand != testCar.Brand || got.Model != testCar.Model {
		t.Fatalf("получен неверный автомобиль: %+v", got)
	}

	// Изменение статуса продажи.
	if err := repo.MarkCarSold(ctx, testCar.ID); err != nil {
		t.Fatalf("MarkCarSold() error: %v", err)
	}

	got, err = repo.GetCarByID(ctx, testCar.ID)
	if err != nil {
		t.Fatalf("GetCarByID() error after sale: %v", err)
	}
	if !got.Sold {
		t.Fatal("ожидалось, что Sold будет true")
	}

	// Проверка отсутствующего автомобиля.
	_, err = repo.GetCarByID(ctx, "unknown-id")
	if err != car.ErrCarNotFound {
		t.Fatalf("ожидалась ErrCarNotFound, получено: %v", err)
	}
}

func BenchmarkGetCarByID(b *testing.B) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	testCar := car.Car{
		ID:    "benchmark-id",
		Brand: "Toyota",
		Model: "Camry",
		Year:  2022,
		Price: 25000,
	}

	if err := repo.CreateCar(ctx, testCar); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := repo.GetCarByID(ctx, testCar.ID); err != nil {
			b.Fatal(err)
		}
	}
}
