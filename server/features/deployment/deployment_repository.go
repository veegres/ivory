package deployment

import (
	"ivory/clients/storage"

	"github.com/google/uuid"
)

// Repository stores custom templates only. The shipped defaults are computed
// from the plugin registries on every request, so they never enter the bucket
// and cannot be edited or deleted by construction.
type Repository struct {
	bucket *storage.DbBucket[Template]
}

func NewRepository(bucket *storage.DbBucket[Template]) *Repository {
	return &Repository{bucket: bucket}
}

func (r *Repository) Get(key uuid.UUID) (Template, error) {
	template, err := r.bucket.Get(key.String())
	return template.withCreation(), err
}

func (r *Repository) List(criteria ListRequest) ([]Template, error) {
	list, err := r.bucket.GetList(func(t Template) bool {
		if criteria.Keeper != nil && t.Keeper != *criteria.Keeper {
			return false
		}
		if criteria.Platform != nil && t.Platform != *criteria.Platform {
			return false
		}
		return true
	}, func(list []Template, i, j int) bool {
		return list[i].CreatedAt < list[j].CreatedAt
	})
	return withCreation(list), err
}

// GetByName looks a name up within one keeper/platform pair, because that is
// the scope the template list - and so the user - sees it in.
func (r *Repository) GetByName(name string, keeper KeeperPlugin, platform PlatformPlugin) (*Template, error) {
	list, err := r.bucket.GetList(func(t Template) bool {
		return t.Name == name && t.Keeper == keeper && t.Platform == platform
	}, nil)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// Create stores a template under the id the service assigned it, mirroring
// Update - the id belongs to the caller so both paths read the same way.
func (r *Repository) Create(key uuid.UUID, template Template) (*Template, error) {
	template.Id = key.String()
	created, err := r.bucket.Create(key.String(), template)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *Repository) Update(key uuid.UUID, template Template) (*Template, error) {
	template.Id = key.String()
	if err := r.bucket.Update(key.String(), template); err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *Repository) Delete(key uuid.UUID) error {
	return r.bucket.Delete(key.String())
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func withCreation(list []Template) []Template {
	for i := range list {
		list[i] = list[i].withCreation()
	}
	return list
}
