package repository

type InitDbRepo interface {
	CreateUserTableIfNotExist() error
	CreateFileTableIfNotExist() error
	InsertSampleData() error
	TruncateAllTables() error
}
