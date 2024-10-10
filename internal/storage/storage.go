package storage

import (
	"errors"
	"golang-island/internal/data"
	"golang.org/x/exp/slog"
	"sync"
)

// ConcurrentMap реализация с блокировкой всей мапы
type ConcurrentMap[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]V
}

func (store *ConcurrentMap[K, V]) Set(key K, value V) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.m[key] = value
}

func (store *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	store.mu.Lock()
	defer store.mu.Unlock()
	val, exists := store.m[key]
	return val, exists
}

func (store *ConcurrentMap[K, V]) Delete(key K) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.m, key)
}

func (store *ConcurrentMap[K, V]) GetAll() map[K]V {
	store.mu.Lock()
	defer store.mu.Unlock()

	// Копируем данные для безопасного возврата
	copyMap := make(map[K]V)
	for key, value := range store.m {
		copyMap[key] = value
	}
	return copyMap
}

type Repository interface {
	PutIfEmpty(id int, data data.Data) bool

	Put(id int, data data.Data)

	ExistsById(id int) bool

	GetById(id int) data.Data
}

type MockRepository struct {
	store *ConcurrentMap[int, data.Data]

	log *slog.Logger
}

func (repo *MockRepository) PutIfEmpty(id int, data data.Data) bool {
	if _, exists := repo.store.Get(id); exists {
		return false // элемент уже существует, не добавляем
	}

	repo.store.Set(id, data)
	return true
}

func (repo *MockRepository) Put(id int, data data.Data) {

	repo.store.Set(id, data)
}

func (repo *MockRepository) ExistsById(id int) bool {

	if _, exists := repo.store.Get(id); exists {
		return true
	}

	return false
}

func (repo *MockRepository) GetById(id int) (data.Data, error) {

	if data, exists := repo.store.Get(id); exists {
		return data, nil
	}

	return data.Data{}, errors.New("not found")
}

func NewRepo(log *slog.Logger) *MockRepository {
	return &MockRepository{
		store: &ConcurrentMap[int, data.Data]{
			m: make(map[int]data.Data),
		},
		log: log,
	}
}
