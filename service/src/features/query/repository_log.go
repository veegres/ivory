package query

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"ivory/src/clients/database"

	"github.com/google/uuid"
)

func (r *Repository) getLog(uuid uuid.UUID) ([]database.QueryFields, error) {
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
		return nil, errors.Join(ErrCannotParseFileCorrupted, errScanner)
	}

	return elements, nil
}

func (r *Repository) hasLog(uuid uuid.UUID) bool {
	return r.queryLogFiles.ExistByName(uuid.String())
}

func (r *Repository) addLog(uuid uuid.UUID, element any) error {
	if !r.queryLogFiles.ExistByName(uuid.String()) {
		_, errCreate := r.queryLogFiles.CreateByName(uuid.String())
		if errCreate != nil {
			return errCreate
		}
	}

	file, err := r.queryLogFiles.OpenByName(uuid.String())
	defer func() { _ = file.Close() }()
	if err != nil {
		return err
	}

	scanBuf := make([]byte, r.maxBufferCapacity)

	counter := bufio.NewScanner(file)
	counter.Buffer(scanBuf, r.maxBufferCapacity)
	lines := 0
	for counter.Scan() {
		lines++
	}
	if errCounter := counter.Err(); errCounter != nil {
		return errors.Join(ErrCannotParseFileCorrupted, errCounter)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

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
		return errors.Join(ErrCannotParseFileCorrupted, errScanner)
	}

	jsonByte, errMarshall := json.Marshal(element)
	if errMarshall != nil {
		return errMarshall
	}
	buf.Write(jsonByte)
	buf.WriteString("\n")

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

func (r *Repository) deleteLog(uuid uuid.UUID) error {
	return r.queryLogFiles.DeleteByName(uuid.String())
}

func (r *Repository) deleteAllLogs() error {
	return r.queryLogFiles.DeleteAll()
}
