package clickhouse

import (
	"context"
	"fmt"
	chclient "ivory/clients/clickhouse"
	"ivory/plugins/database"
	"reflect"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func (a *Adapter) GetMany(ctx database.Context, query string, queryParams []any) ([]string, error) {
	conn, _, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer chclient.Close(conn)

	requestCtx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	rows, errQuery := conn.Query(requestCtx, query, queryParams...)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	columnTypes := rows.ColumnTypes()
	results := make([]string, 0)
	for rows.Next() {
		dest := scanDestinations(columnTypes)
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		if len(dest) > 0 {
			results = append(results, fmt.Sprintf("%v", reflect.ValueOf(dest[0]).Elem().Interface()))
		}
	}
	return results, rows.Err()
}

func (a *Adapter) GetOne(ctx database.Context, query string) (any, error) {
	conn, _, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer chclient.Close(conn)

	requestCtx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	rows, errQuery := conn.Query(requestCtx, query)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, rows.Err()
	}

	columnTypes := rows.ColumnTypes()
	if len(columnTypes) == 0 {
		return nil, nil
	}
	dest := scanDestinations(columnTypes)
	if err := rows.Scan(dest...); err != nil {
		return nil, err
	}
	return reflect.ValueOf(dest[0]).Elem().Interface(), nil
}

func (a *Adapter) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	finalQuery := query
	var params []any
	resultOptions := options
	if options != nil {
		newQuery, newLimit, errNormalize := a.normalizeQuery(query, options.Trim, options.Limit)
		if errNormalize != nil {
			return nil, errNormalize
		}
		finalQuery = newQuery
		params = options.Params
		clone := *options
		clone.Limit = newLimit
		resultOptions = &clone
	}

	conn, url, errConnect := a.connect(ctx)
	if errConnect != nil {
		return nil, errConnect
	}
	defer chclient.Close(conn)

	requestCtx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	rows, errQuery := conn.Query(requestCtx, finalQuery, params...)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()

	fields, resultRows, errScan := scanRows(rows)
	if errScan != nil {
		return nil, errScan
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &database.QueryFields{
		Fields:    fields,
		Rows:      resultRows,
		Url:       url,
		StartTime: startTime,
		EndTime:   time.Now().UnixMilli(),
		Options:   resultOptions,
	}, nil
}

// scanDestinations builds one addressable pointer per column, each of the
// exact Go type the driver reports it will decode that column into, so Scan
// never has to guess a type for a generic destination.
func scanDestinations(columnTypes []driver.ColumnType) []any {
	dest := make([]any, len(columnTypes))
	for i, ct := range columnTypes {
		dest[i] = reflect.New(ct.ScanType()).Interface()
	}
	return dest
}

func scanRows(rows driver.Rows) ([]database.QueryField, [][]any, error) {
	columnTypes := rows.ColumnTypes()
	names := rows.Columns()

	fields := make([]database.QueryField, len(names))
	for i, name := range names {
		fields[i] = database.QueryField{Name: name, DataType: columnTypes[i].DatabaseTypeName(), DataTypeOID: 0}
	}

	resultRows := make([][]any, 0)
	for rows.Next() {
		dest := scanDestinations(columnTypes)
		if err := rows.Scan(dest...); err != nil {
			return nil, nil, err
		}
		row := make([]any, len(dest))
		for i, d := range dest {
			row[i] = reflect.ValueOf(d).Elem().Interface()
		}
		resultRows = append(resultRows, row)
	}
	return fields, resultRows, nil
}
