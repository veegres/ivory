package query

type DatabaseClient interface {
	GetApplicationName(session string) string
	GetMany(ctx QueryContext, query string, queryParams []any) ([]string, error)
	GetOne(ctx QueryContext, query string) (any, error)
	GetFields(ctx QueryContext, query string, options *QueryOptions) (*QueryFields, error)
	Cancel(ctx QueryContext, pid int) error
	Terminate(ctx QueryContext, pid int) error
}
