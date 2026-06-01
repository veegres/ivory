package utils

import "errors"

var ErrNotImplemented = errors.New("adapter is not implemented")

type Registry[K comparable, V any] struct {
	items map[K]V
}

func NewRegistry[K comparable, V any]() *Registry[K, V] {
	return &Registry[K, V]{
		items: make(map[K]V),
	}
}

func (r *Registry[K, V]) Register(name K, item V) {
	r.items[name] = item
}

func (r *Registry[K, V]) Get(name K) (V, error) {
	if item, ok := r.items[name]; ok {
		return item, nil
	}

	var zero V
	return zero, ErrNotImplemented
}

func (r *Registry[K, V]) All() map[K]V {
	return r.items
}
