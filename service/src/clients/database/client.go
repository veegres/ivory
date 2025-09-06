package database

type Client interface {
	GetApplicationName(session string) string
	GetMany(ctx Context, query string, queryParams []any) ([]string, error)
	GetOne(ctx Context, query string) (any, error)
	GetFields(ctx Context, query string, options *QueryOptions) (*QueryFields, error)
	Cancel(ctx Context, pid int) error
	Terminate(ctx Context, pid int) error
}
