package model

type DbConnection struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Target struct {
	DbName        string `json:"dbName"`
	Schema        string `json:"schema"`
	Table         string `json:"table"`
	ExcludeSchema string `json:"excludeSchema"`
	ExcludeTable  string `json:"excludeTable"`
	Ratio         int8   `json:"ratio"`
}

type PgCompactTable struct {
	Connection DbConnection `json:"connection"`
	Target     *Target      `json:"target"`
}
