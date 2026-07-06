package postgres

import (
	"ivory/plugins/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (a *Adapter) GetMany(ctx database.Context, query string, queryParams []any) ([]string, error) {
	values := make([]string, 0)
	errReq := a.sendRequest(ctx, query, queryParams, func(rows pgx.Rows, _ *pgtype.Map, _ string) error {
		for rows.Next() {
			var value string
			err := rows.Scan(&value)
			if err != nil {
				return err
			}
			values = append(values, value)
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}
	return values, nil
}

func (a *Adapter) GetOne(ctx database.Context, query string) (any, error) {
	var value any
	errReq := a.sendRequest(ctx, query, nil, func(rows pgx.Rows, _ *pgtype.Map, _ string) error {
		if rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return err
			}
			value = values[0]
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}
	return value, nil
}

func (a *Adapter) GetFields(ctx database.Context, query string, options *database.QueryOptions) (*database.QueryFields, error) {
	startTime := time.Now().UnixMilli()

	fields := make([]database.QueryField, 0)
	rowList := make([][]any, 0)
	url := "-"

	// NOTE: we need this object ot avoid fatal errors when the option variable is `nil`
	tmpOptions := database.QueryOptions{}
	if options != nil {
		tmpOptions.Limit = options.Limit
		tmpOptions.Trim = options.Trim
		tmpOptions.Params = options.Params
	}

	normQuery, normLimit, errNormQuery := a.normalizeQuery(query, tmpOptions.Trim, tmpOptions.Limit)
	if errNormQuery != nil {
		return nil, errNormQuery
	}

	errReq := a.sendRequest(ctx, normQuery, tmpOptions.Params, func(rows pgx.Rows, typeMap *pgtype.Map, connUrl string) error {
		url = connUrl
		for _, field := range rows.FieldDescriptions() {
			dataType, ok := typeMap.TypeForOID(field.DataTypeOID)
			if !ok {
				fields = append(fields, database.QueryField{Name: field.Name, DataType: "unknown", DataTypeOID: field.DataTypeOID})
			} else {
				fields = append(fields, database.QueryField{Name: field.Name, DataType: dataType.Name, DataTypeOID: field.DataTypeOID})
			}
		}
		for rows.Next() {
			val, err := rows.Values()
			if err != nil {
				return err
			}
			rowList = append(rowList, val)
		}
		return nil
	})
	if errReq != nil {
		return nil, errReq
	}

	if options != nil {
		options.Limit = normLimit
	}
	endTime := time.Now().UnixMilli()
	res := &database.QueryFields{
		Fields:    fields,
		Rows:      rowList,
		Url:       url,
		StartTime: startTime,
		EndTime:   endTime,
		Options:   options,
	}
	return res, nil
}
