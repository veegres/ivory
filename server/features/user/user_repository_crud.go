package user

import "sort"

func (r *Repository) Get(username string) (User, error) {
	return r.bucket.Get(username)
}

func (r *Repository) List() ([]User, error) {
	return r.bucket.GetList(nil, func(list []User, i, j int) bool {
		return list[i].Username < list[j].Username
	})
}

func (r *Repository) Create(user User) (User, error) {
	return r.bucket.Create(user.Username, user)
}

func (r *Repository) Update(user User) error {
	return r.bucket.Update(user.Username, user)
}

func (r *Repository) Delete(username string) error {
	return r.bucket.Delete(username)
}

// DeleteIf deletes a user in the same transaction that checks the rule about
// who is left, so two deletes at once cannot both believe they are not the last.
func (r *Repository) DeleteIf(username string, check func(users map[string]User) error) error {
	return r.bucket.DeleteIf(username, check)
}

func (r *Repository) DeleteAll() error {
	return r.bucket.DeleteAll()
}

func (r *Repository) Superusers() ([]string, error) {
	users, err := r.bucket.GetList(func(el User) bool { return el.Superuser }, nil)
	if err != nil {
		return nil, err
	}
	superusers := make([]string, 0, len(users))
	for _, u := range users {
		superusers = append(superusers, u.Username)
	}
	sort.Strings(superusers)
	return superusers, nil
}
