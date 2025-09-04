package query

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"ivory/src/clients/database/bolt"
	"ivory/src/clients/database/files"
	"time"

	"github.com/google/uuid"
)

type QueryRepository struct {
	bucket            *bolt.Bucket[Query]
	queryLogFiles     *files.FileGateway
	maxBufferCapacity int
	maxLogElements    int
}

func NewQueryRepository(
	bucket *bolt.Bucket[Query],
	queryLogFiles *files.FileGateway,
) *QueryRepository {
	return &QueryRepository{
		bucket:            bucket,
		queryLogFiles:     queryLogFiles,
		maxBufferCapacity: 1024 * 1024,
		maxLogElements:    10,
	}
}

func (r *QueryRepository) GetLog(uuid uuid.UUID) ([]QueryFields, error) {
	file, err := r.queryLogFiles.OpenByName(uuid.String())
	if err != nil {
		return []QueryFields{}, nil
	}

	var elements []QueryFields
	scanner := bufio.NewScanner(file)
	buf := make([]byte, r.maxBufferCapacity)
	scanner.Buffer(buf, r.maxBufferCapacity)
	for scanner.Scan() {
		b := scanner.Bytes()
		if len(b) != 0 {
			var obj QueryFields
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

func (r *QueryRepository) DeleteLog(uuid uuid.UUID) error {
	return r.queryLogFiles.DeleteByName(uuid.String())
}

func (r *QueryRepository) AddLog(uuid uuid.UUID, element any) error {
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
	keyString := key.String()
	var errFile error
	if r.queryLogFiles.ExistByName(keyString) {
		errFile = r.queryLogFiles.DeleteByName(keyString)
	}
	errBucket := r.bucket.Delete(keyString)
	return errors.Join(errFile, errBucket)
}

func (r *QueryRepository) DeleteAll() error {
	errFile := r.queryLogFiles.DeleteAll()
	errBucket := r.bucket.DeleteAll()
	return errors.Join(errFile, errBucket)
}

func (r *QueryRepository) sortAscByCreatedAt(list []Query, i, j int) bool {
	return list[i].CreatedAt < list[j].CreatedAt
}
