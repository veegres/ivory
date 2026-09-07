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
		if plugin != nil && !r.matchesPlugin(query.Plugin, *plugin) {
			return false
		}
		return true
	}, r.sortAscByCreatedAt)
}

// SystemQueryNamesForPlugin returns the names of the system queries already
// stored for a plugin. A name is the stable identity of a system query - it is
// the one field Update refuses to change - so it is what seeding compares
// against to add the queries a new version ships without duplicating, or
// overwriting an edit made to, the ones already there.
func (r *Repository) SystemQueryNamesForPlugin(plugin DbPlugin) (map[string]bool, error) {
	list, err := r.bucket.GetList(func(query Response) bool {
		return query.Creation == System && r.matchesPlugin(query.Plugin, plugin)
	}, nil)
	if err != nil {
		return nil, err
	}
	names := make(map[string]bool, len(list))
	for _, query := range list {
		names[query.Name] = true
	}
	return names, nil
}

// matchesPlugin treats records stored before the plugin field existed
// (empty plugin) as postgres queries.
func (r *Repository) matchesPlugin(stored DbPlugin, requested DbPlugin) bool {
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
