package query

import (
	"ivory/plugins/database"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) Get(uuid uuid.UUID) (Response, error) {
	return r.bucket.Get(uuid.String())
}

func (r *Repository) List() ([]Response, error) {
	return r.ListByFilter(nil, nil)
}

func (r *Repository) ListByType(queryType Type) ([]Response, error) {
	return r.ListByFilter(&queryType, nil)
}

// ListByFilter filters by type and plugin when they are provided; nil means
// no filtering for that dimension.
func (r *Repository) ListByFilter(queryType *Type, plugin *DbPlugin) ([]Response, error) {
	return r.bucket.GetList(func(query Response) bool {
		if queryType != nil && query.Type != *queryType {
			return false
		}
		if plugin != nil && !matchesPlugin(query.Plugin, *plugin) {
			return false
		}
		return true
	}, r.sortAscByCreatedAt)
}

func (r *Repository) HasSystemQueriesForPlugin(plugin DbPlugin) (bool, error) {
	list, err := r.bucket.GetList(func(query Response) bool {
		return query.Creation == System && matchesPlugin(query.Plugin, plugin)
	}, nil)
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}

// matchesPlugin treats records stored before the plugin field existed
// (empty plugin) as postgres queries.
func matchesPlugin(stored DbPlugin, requested DbPlugin) bool {
	if stored == "" {
		stored = database.POSTGRES
	}
	return stored == requested
}

func (r *Repository) Create(query Response) (*uuid.UUID, *Response, error) {
	key := uuid.New()
	query.Id = key
	query.CreatedAt = time.Now().UnixNano()
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *Repository) Update(key uuid.UUID, query Response) (*uuid.UUID, *Response, error) {
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *Repository) Delete(key uuid.UUID) error {
	return r.bucket.Delete(key.String())
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *Repository) sortAscByCreatedAt(list []Response, i, j int) bool {
	return list[i].CreatedAt < list[j].CreatedAt
}
