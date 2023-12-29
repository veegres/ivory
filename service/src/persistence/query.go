package persistence

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"ivory/src/config"
	. "ivory/src/model"
	"time"
)

type QueryRepository struct {
	bucket            *config.Bucket[Query]
	queryHistoryFiles *config.FileGateway
}

func NewQueryRepository(
	bucket *config.Bucket[Query],
	queryHistoryFiles *config.FileGateway,
) *QueryRepository {
	return &QueryRepository{
		bucket:            bucket,
		queryHistoryFiles: queryHistoryFiles,
	}
}

func (r *QueryRepository) GetHistory(uuid uuid.UUID) ([]QueryFields, error) {
	file, err := r.queryHistoryFiles.Open(uuid.String())
	if err != nil {
		return []QueryFields{}, nil
	}

	var elements []QueryFields
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		if len(bytes) != 0 {
			var obj QueryFields
			errUnmarshal := json.Unmarshal(scanner.Bytes(), &obj)
			if errUnmarshal == nil {
				elements = append(elements, obj)
			}
		}
	}

	if errScanner := scanner.Err(); errScanner != nil {
		return nil, errScanner
	}

	return elements, nil
}

func (r *QueryRepository) DeleteHistory(uuid uuid.UUID) error {
	return r.queryHistoryFiles.Delete(uuid.String())
}

func (r *QueryRepository) AddHistory(uuid uuid.UUID, element any) error {
	if !r.queryHistoryFiles.Exist(uuid.String()) {
		_, errCreate := r.queryHistoryFiles.Create(uuid.String())
		if errCreate != nil {
			return errCreate
		}
	}

	file, err := r.queryHistoryFiles.Open(uuid.String())
	defer func() { _ = file.Close() }()
	if err != nil {
		return err
	}

	jsonByte, errMarshall := json.Marshal(element)
	if errMarshall != nil {
		return errMarshall
	}

	_, errWrite := file.Write(jsonByte)
	if errWrite != nil {
		return errWrite
	}
	_, errEnter := file.WriteString("\n")
	if errEnter != nil {
		return errEnter
	}

	return nil
}

func (r *QueryRepository) Get(uuid uuid.UUID) (Query, error) {
	return r.bucket.Get(uuid.String())
}

func (r *QueryRepository) List() ([]Query, error) {
	return r.bucket.GetList(nil, r.sortAscByCreatedAt)
}

func (r *QueryRepository) ListByType(queryType QueryType) ([]Query, error) {
	return r.bucket.GetList(func(cert Query) bool {
		return cert.Type == queryType
	}, r.sortAscByCreatedAt)
}

func (r *QueryRepository) Create(query Query) (*uuid.UUID, *Query, error) {
	key := uuid.New()
	query.Id = key
	query.CreatedAt = time.Now().UnixNano()
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *QueryRepository) Update(key uuid.UUID, query Query) (*uuid.UUID, *Query, error) {
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *QueryRepository) Delete(key uuid.UUID) error {
	errFile := r.queryHistoryFiles.Delete(key.String())
	errBucket := r.bucket.Delete(key.String())
	return errors.Join(errFile, errBucket)
}

func (r *QueryRepository) DeleteAll() error {
	errFile := r.queryHistoryFiles.DeleteAll()
	errBucket := r.bucket.DeleteAll()
	return errors.Join(errFile, errBucket)
}

func (r *QueryRepository) sortAscByCreatedAt(list []Query, i, j int) bool {
	return list[i].CreatedAt < list[j].CreatedAt
}
