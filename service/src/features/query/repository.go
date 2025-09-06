package query

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"ivory/src/clients/database"
	"ivory/src/storage/bolt"
	"ivory/src/storage/files"
	"time"

	"github.com/google/uuid"
)

type LogRepository struct {
	queryLogFiles     *files.FileGateway
	maxBufferCapacity int
	maxLogElements    int
}

func NewLogRepository(
	queryLogFiles *files.FileGateway,
) *LogRepository {
	return &LogRepository{
		queryLogFiles:     queryLogFiles,
		maxBufferCapacity: 1024 * 1024,
		maxLogElements:    10,
	}
}

func (r *LogRepository) Get(uuid uuid.UUID) ([]database.QueryFields, error) {
	file, err := r.queryLogFiles.OpenByName(uuid.String())
	if err != nil {
		return []database.QueryFields{}, nil
	}

	var elements []database.QueryFields
	scanner := bufio.NewScanner(file)
	buf := make([]byte, r.maxBufferCapacity)
	scanner.Buffer(buf, r.maxBufferCapacity)
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) != 0 {
			var obj database.QueryFields
			errUnmarshal := json.Unmarshal(scanner.Bytes(), &obj)
			if errUnmarshal == nil {
				elements = append(elements, obj)
			}
		}
	}

	if errScanner := scanner.Err(); errScanner != nil {
		return nil, errors.Join(errors.New("cannot parse file, it is corrupted"), errScanner)
	}

	return elements, nil
}

func (r *LogRepository) Exist(uuid uuid.UUID) bool {
	return r.queryLogFiles.ExistByName(uuid.String())
}

func (r *LogRepository) Add(uuid uuid.UUID, element any) error {
	if !r.queryLogFiles.ExistByName(uuid.String()) {
		_, errCreate := r.queryLogFiles.CreateByName(uuid.String())
		if errCreate != nil {
			return errCreate
		}
	}

	// open file
	file, err := r.queryLogFiles.OpenByName(uuid.String())
	defer func() { _ = file.Close() }()
	if err != nil {
		return err
	}

	scanBuf := make([]byte, r.maxBufferCapacity)

	// count number of lines in is the file
	counter := bufio.NewScanner(file)
	counter.Buffer(scanBuf, r.maxBufferCapacity)
	lines := 0
	for counter.Scan() {
		lines++
	}
	if errCounter := counter.Err(); errCounter != nil {
		return errors.Join(errors.New("cannot parse file, it is corrupted"), errCounter)
	}

	// reset cursor in the file for new scanner
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// removing old rows to fit max log elements
	buf := bytes.NewBuffer([]byte{})
	scanner := bufio.NewScanner(file)
	scanner.Buffer(scanBuf, r.maxBufferCapacity)
	for scanner.Scan() {
		if lines < r.maxLogElements {
			buf.Write(scanner.Bytes())
			buf.WriteString("\n")
		} else {
			lines--
		}
	}
	if errScanner := scanner.Err(); errScanner != nil {
		return errors.Join(errors.New("cannot parse file, it is corrupted"), errScanner)
	}

	// parse and add new element to the buffer
	jsonByte, errMarshall := json.Marshal(element)
	if errMarshall != nil {
		return errMarshall
	}
	buf.Write(jsonByte)
	buf.WriteString("\n")

	// clean and rewrite a file
	errTruncate := file.Truncate(0)
	if errTruncate != nil {
		return errTruncate
	}
	_, errWrite := file.Write(buf.Bytes())
	if errWrite != nil {
		return errWrite
	}

	return nil
}

func (r *LogRepository) Delete(uuid uuid.UUID) error {
	return r.queryLogFiles.DeleteByName(uuid.String())
}

func (r *LogRepository) DeleteAll() error {
	return r.queryLogFiles.DeleteAll()
}

type Repository struct {
	bucket        *bolt.Bucket[Query]
	LogRepository *LogRepository
}

func NewRepository(
	bucket *bolt.Bucket[Query],
	logRepository *LogRepository,
) *Repository {
	return &Repository{
		bucket:        bucket,
		LogRepository: logRepository,
	}
}

func (r *Repository) Get(uuid uuid.UUID) (Query, error) {
	return r.bucket.Get(uuid.String())
}

func (r *Repository) List() ([]Query, error) {
	return r.bucket.GetList(nil, r.sortAscByCreatedAt)
}

func (r *Repository) ListByType(queryType database.QueryType) ([]Query, error) {
	return r.bucket.GetList(func(cert Query) bool {
		return cert.Type == queryType
	}, r.sortAscByCreatedAt)
}

func (r *Repository) Create(query Query) (*uuid.UUID, *Query, error) {
	key := uuid.New()
	query.Id = key
	query.CreatedAt = time.Now().UnixNano()
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *Repository) Update(key uuid.UUID, query Query) (*uuid.UUID, *Query, error) {
	err := r.bucket.Update(key.String(), query)
	return &key, &query, err
}

func (r *Repository) Delete(key uuid.UUID) error {
	var errFile error
	if r.LogRepository.Exist(key) {
		errFile = r.LogRepository.Delete(key)
	}
	errBucket := r.bucket.Delete(key.String())
	return errors.Join(errFile, errBucket)
}

func (r *Repository) DeleteAll() error {
	errFile := r.LogRepository.DeleteAll()
	errBucket := r.bucket.DeleteAll()
	return errors.Join(errFile, errBucket)
}

func (r *Repository) sortAscByCreatedAt(list []Query, i, j int) bool {
	return list[i].CreatedAt < list[j].CreatedAt
}
