package user

func (r *Repository) LinkGet(id string) (Link, error) {
	return r.linkBucket.Get(id)
}

func (r *Repository) LinkMap() (map[string]Link, error) {
	return r.linkBucket.GetMap(nil)
}

func (r *Repository) LinkCreate(id string, link Link) (Link, error) {
	return r.linkBucket.Create(id, link)
}

func (r *Repository) LinkDelete(id string) error {
	return r.linkBucket.Delete(id)
}

func (r *Repository) LinkDeleteAll() error {
	return r.linkBucket.DeleteAll()
}
