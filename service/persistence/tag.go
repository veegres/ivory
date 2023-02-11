package persistence

type TagRepository struct {
	common common
	bucket []byte
}

func (t TagRepository) List() ([]string, error) {
	return GetKeyList(t.bucket)
}

func (t TagRepository) Get(tag string) ([]string, error) {
	return Get[[]string](t.bucket, tag)
}

func (t TagRepository) GetMap() (map[string][]string, error) {
	return GetMap[[]string](t.bucket)
}

func (t TagRepository) UpdateCluster(cluster string, tags []string) error {
	tagMap, err := t.GetMap()
	if err != nil {
		return err
	}

	// NOTE: remove cluster from all tags
	for k, v := range tagMap {
		tmp := make([]string, 0)
		for _, c := range v {
			if c != cluster {
				tmp = append(tmp, c)
			}
		}

		if len(tmp) == 0 {
			tagMap[k] = nil
		} else {
			tagMap[k] = tmp
		}
	}

	// NOTE: add cluster to tags
	for _, v := range tags {
		tagMap[v] = append(tagMap[v], cluster)
	}

	// NOTE: update tags
	for k, v := range tagMap {
		if v != nil {
			err := Update[[]string](t.bucket, k, v)
			if err != nil {
				return err
			}
		} else {
			err := Delete(t.bucket, k)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t TagRepository) Delete(tag string) error {
	return Delete(t.bucket, tag)
}

func (t TagRepository) DeleteAll() error {
	return DeleteAll(t.bucket)
}
